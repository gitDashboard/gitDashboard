package admin

import (
	"github.com/gitDashboard/client/v1/admin/request"
	"github.com/gitDashboard/client/v1/admin/response"
	basicResponse "github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/gitDashboard/gitDashboard/app/repoManager"
	"github.com/revel/revel"
	git "gopkg.in/libgit2/git2go.v22"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type AdminRepo struct {
	controllers.AdminController
}

func (ctrl *AdminRepo) CreateFolder() revel.Result {
	var req request.CreateFolderRequest
	var resp response.CreateFolderResponse
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		return ctrl.RenderError(err)
	}
	fullPath := controllers.GitBasePath() + "/" + req.Path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.Mkdir(fullPath, 0770)
		if err != nil {
			controllers.ErrorResp(&resp, basicResponse.FatalError, err)
		} else {
			resp.Success = true
		}
	} else {
		resp.Success = false
		resp.Error = basicResponse.AlreadyExistError
	}
	return ctrl.RenderJson(resp)
}

func ConfigRepo(repo *git.Repository) error {
	repoConfig, _ := repo.Config()
	repoConfig.SetString("core.sharedRepository", "group")
	//create update hook as symbolic link
	hooks := []string{"update", "post-update"}
	for _, hook := range hooks {
		hookPath := revel.Config.StringDefault("gd.hookFolder", "") + "/" + hook
		repoHookPath := repo.Path() + "/hooks/" + hook
		if _, existHookErr := os.Stat(hookPath); existHookErr == nil {
			if _, errHookAlreadyExist := os.Stat(repoHookPath); os.IsNotExist(errHookAlreadyExist) {
				revel.INFO.Printf("Creatin link from %s to %s \n", hookPath, repoHookPath)
				return os.Symlink(hookPath, repo.Path()+"/hooks/"+hook)
			}
		}
	}
	return nil
}

func UpdateRepoDescription(repoPath, description string) {
	ioutil.WriteFile(repoPath+"/description", []byte(description), 0770)
}

func (ctrl *AdminRepo) UpdateDescription(repoId uint) revel.Result {
	var resp basicResponse.BasicResponse
	var dbRepo models.Repo
	var description string
	ctrl.Params.Bind(&description, "description")
	if description != "" {
		ctrl.Tx.First(&dbRepo, repoId)
		if dbRepo.ID != repoId {
			controllers.ErrorResp(&resp, basicResponse.NoRepositoryFoundError, nil)
			return ctrl.RenderJson(resp)
		}
		UpdateRepoDescription(dbRepo.Path, description)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminRepo) CreateRepo() revel.Result {
	var req request.CreateRepoRequest
	var resp response.CreateRepoResponse
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		return ctrl.RenderError(err)
	}
	fullPath := controllers.GitBasePath() + "/"
	if !strings.HasSuffix(req.Path, ".git") {
		fullPath = fullPath + req.Path + ".git"
	} else {
		fullPath = fullPath + req.Path
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		repo, err := git.InitRepository(fullPath, true)
		if err != nil {
			controllers.ErrorResp(&resp, basicResponse.FatalError, err)
			return ctrl.RenderJson(resp)
		}
		defer repo.Free()
		if req.Description != "" {
			UpdateRepoDescription(fullPath, req.Description)
		}
		err = ConfigRepo(repo)
		if err != nil {
			panic(err)
		}
		dbRepo := &models.Repo{Path: controllers.CleanSlashes(fullPath)}
		ctrl.Tx.Create(dbRepo)
		resp.Success = true
	} else {
		controllers.ErrorResp(&resp, basicResponse.AlreadyExistError, nil)
	}
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminRepo) InitExistingRepo() revel.Result {
	var resp response.CreateFolderResponse
	var repoPath string
	ctrl.Params.Bind(&repoPath, "path")
	repoPath = controllers.CleanSlashes(controllers.GitBasePath() + "/" + repoPath)
	//check if path is a git repository
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		controllers.ErrorResp(&resp, basicResponse.NoRepositoryFoundError, nil)
		return ctrl.RenderJson(resp)
	}
	err = ConfigRepo(repo)
	if err != nil {
		controllers.ErrorResp(&resp, basicResponse.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	dbRepo := models.Repo{Path: repoPath}
	db := ctrl.Tx.Create(&dbRepo)
	if db.Error != nil {
		controllers.ErrorResp(&resp, basicResponse.DbError, db.Error)
		return ctrl.RenderJson(resp)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminRepo) Permissions(repoId uint) revel.Result {
	var resp response.GetPermissionsResponse
	var dbRepo models.Repo
	ctrl.Tx.First(&dbRepo, repoId)
	if len(ctrl.Tx.GetErrors()) > 0 {
		controllers.ErrorResp(&resp, basicResponse.FatalError, nil)
		return ctrl.RenderJson(resp)
	}
	if dbRepo.ID != uint(repoId) {
		controllers.ErrorResp(&resp, basicResponse.NoRepositoryFoundError, nil)
		return ctrl.RenderJson(resp)
	}
	ctrl.Tx.Order("position").Where("repo_id=?", dbRepo.ID).Find(&dbRepo.Permissions)

	resp.Permissions = make([]response.RepoPermission, len(dbRepo.Permissions), len(dbRepo.Permissions))
	for i, perm := range dbRepo.Permissions {
		var repoPerm response.RepoPermission
		repoPerm.Ref = perm.Branch
		repoPerm.Position = perm.Position
		repoPerm.Granted = perm.Granted
		if perm.UserID.Valid {
			repoPerm.UserID = perm.UserID.Int64
			//search user
			var dbUser models.User
			ctrl.Tx.First(&dbUser, perm.UserID)
			repoPerm.UserName = dbUser.Username
		}
		if perm.GroupID.Valid {
			repoPerm.GroupID = perm.GroupID.Int64
			//search group
			var dbGroup models.Group
			ctrl.Tx.First(&dbGroup, perm.GroupID)
			repoPerm.GroupName = dbGroup.Name
		}
		if perm.Type != "" {
			repoPerm.Types = strings.Split(perm.Type, ",")
		}
		resp.Permissions[i] = repoPerm
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminRepo) UpdatePermissions(repoId uint) revel.Result {
	var req request.UpdatePermissionsRequest
	var resp response.UpdatePermissionsResponse
	var dbRepo models.Repo
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ctrl.RenderError(err)
	}

	revel.INFO.Printf("UpdatePermissions req:%+v\n", req)

	db := ctrl.Tx.First(&dbRepo, repoId)
	if len(db.GetErrors()) > 0 {
		controllers.ErrorResp(&resp, basicResponse.FatalError, db.GetErrors()[0])
		return ctrl.RenderJson(resp)
	}
	if dbRepo.ID != uint(repoId) {
		controllers.ErrorResp(&resp, basicResponse.NoRepositoryFoundError, nil)
		return ctrl.RenderJson(resp)
	}

	//remove all old permissions
	ctrl.Tx.Where("repo_id = ?", repoId).Delete(models.Permission{})
	//insert all new permissions
	for _, newPerm := range req.Permissions {
		dbPermission := models.Permission{
			RepoID:   repoId,
			Branch:   newPerm.Ref,
			Granted:  newPerm.Granted,
			Position: newPerm.Position,
		}
		for _, permType := range newPerm.Types {
			if permType != "" {
				dbPermission.Type = dbPermission.Type + "," + permType
			}
		}
		dbPermission.Type = dbPermission.Type[1:]
		dbPermission.UserID.Scan(newPerm.UserID)
		dbPermission.GroupID.Scan(newPerm.GroupID)
		db := ctrl.Tx.Create(&dbPermission)
		if len(db.GetErrors()) > 0 {
			controllers.ErrorResp(&resp, basicResponse.FatalError, db.GetErrors()[0])
			return ctrl.RenderJson(resp)
		}
	}

	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminRepo) Move(repoId uint) revel.Result {
	var req request.RepoMoveRequest
	var resp basicResponse.BasicResponse
	var dbRepo models.Repo
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ctrl.RenderError(err)
	}

	db := ctrl.Tx.First(&dbRepo, repoId)
	if db.Error != nil {
		controllers.ErrorResp(&resp, basicResponse.FatalError, db.Error)
		return ctrl.RenderJson(resp)
	}
	if dbRepo.ID != uint(repoId) {
		controllers.ErrorResp(&resp, basicResponse.NoRepositoryFoundError, nil)
		return ctrl.RenderJson(resp)
	}

	//locking db on new transaction
	toUnlock := false
	if !dbRepo.Locked {
		revel.INFO.Println("Locking repo")
		lockDbTx := ctrl.NewTransaction()
		defer ctrl.RollbackTransaction(lockDbTx)
		lockDbTx.Table("repos").Where("id=?", repoId).Update("locked", true)
		ctrl.CommitTransaction(lockDbTx)
		toUnlock = true
		revel.INFO.Println("repo Locked")
	}
	err = ctrl.waitForOperations(dbRepo.ID)
	if err != nil {
		controllers.ErrorResp(&resp, basicResponse.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	//extract reponame
	repoStat, err := os.Stat(dbRepo.Path)
	if err != nil {
		controllers.ErrorResp(&resp, basicResponse.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	var newRepoPath string
	if req.DestName == "" {
		//move
		newRepoPath = controllers.CleanSlashes(controllers.GitBasePath() + "/" + req.DestPath + "/" + repoStat.Name())
	} else {
		//rename
		if req.DestPath == "" {
			req.DestPath = strings.Replace(path.Dir(dbRepo.Path)+"/", controllers.GitBasePath(), "", 1)
		}
		newRepoPath = controllers.CleanSlashes(controllers.GitBasePath() + "/" + req.DestPath + "/" + req.DestName)
	}
	db = ctrl.Tx.Table("repos").Where("id=?", repoId).Update("path", newRepoPath)

	if db.Error != nil {
		controllers.ErrorResp(&resp, basicResponse.FatalError, db.Error)
		return ctrl.RenderJson(resp)
	}
	err = os.Rename(dbRepo.Path, newRepoPath)
	if err != nil {
		go ctrl.unlockRepo(repoId)
		panic(err) //to rollback the transaction
	}
	//unlock repo
	if toUnlock {
		go ctrl.unlockRepo(repoId)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

/** wait until no running operation on the repo */
func (ctrl *AdminRepo) waitForOperations(repoId uint) error {
	hasOp := true
	var err error
	for hasOp {
		hasOp, err = repoManager.HasOperationInProgress(ctrl.Tx, repoId)
		if err != nil {
			return err
		}
		if hasOp {
			time.Sleep(500 * time.Millisecond)
		}
	}
	return nil
}

func (ctrl *AdminRepo) unlockRepo(repoId uint) {
	revel.INFO.Println("UnLocking repo")
	unLockDbTx := ctrl.NewTransaction()
	defer ctrl.RollbackTransaction(unLockDbTx)
	unLockDbTx.Table("repos").Where("id=?", repoId).Update("locked", false)
	ctrl.CommitTransaction(unLockDbTx)
	revel.INFO.Println("repo Unlocked")
}

func (ctrl *AdminRepo) Lock(lock bool, repoId uint) revel.Result {
	var resp basicResponse.BasicResponse
	var dbRepo models.Repo
	db := ctrl.Tx.First(&dbRepo, repoId)
	if db.Error != nil {
		controllers.ErrorResp(&resp, basicResponse.FatalError, db.Error)
		return ctrl.RenderJson(resp)
	}
	if dbRepo.ID != uint(repoId) {
		controllers.ErrorResp(&resp, basicResponse.NoRepositoryFoundError, nil)
		return ctrl.RenderJson(resp)
	}
	db = ctrl.Tx.Model(dbRepo).Update("locked", lock)
	if db.Error != nil {
		controllers.ErrorResp(&resp, basicResponse.FatalError, db.Error)
		return ctrl.RenderJson(resp)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}
