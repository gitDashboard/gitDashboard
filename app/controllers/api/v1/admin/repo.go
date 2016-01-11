package admin

import (
	"github.com/gitDashboard/client/v1/admin/request"
	"github.com/gitDashboard/client/v1/admin/response"
	basicResponse "github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
	git "gopkg.in/libgit2/git2go.v22"
	"io/ioutil"
	"os"
	"strings"
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
			resp.Success = false
			resp.Error = basicResponse.FatalError
			resp.Error.Message = resp.Error.Message + err.Error()
		} else {
			resp.Success = true
		}
	} else {
		resp.Success = false
		resp.Error = basicResponse.AlreadyExistError
	}
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
			resp.Success = false
			resp.Error = basicResponse.FatalError
			resp.Error.Message = resp.Error.Message + err.Error()
			return ctrl.RenderJson(resp)
		}
		defer repo.Free()
		repoConfig, _ := repo.Config()
		repoConfig.SetString("core.sharedRepository", "group")
		if req.Description != "" {
			ioutil.WriteFile(fullPath+"/description", []byte(req.Description), 0770)
		}
		//create update hook as symbolic link
		updateHookPath := revel.Config.StringDefault("gd.hookFolder", "") + "/update"
		if _, existHookErr := os.Stat(updateHookPath); existHookErr == nil {
			revel.INFO.Printf("Creatin link from %s to %s \n", updateHookPath, fullPath+"/hooks/update")
			err = os.Symlink(updateHookPath, fullPath+"/hooks/update")
			if err != nil {
				revel.ERROR.Println(err.Error())
			}
		}
		dbRepo := &models.Repo{Path: controllers.CleanSlashes(fullPath)}
		ctrl.Tx.Create(dbRepo)
		resp.Success = true
	} else {
		resp.Success = false
		resp.Error = basicResponse.AlreadyExistError
	}
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminRepo) Permissions(repoId uint) revel.Result {
	var resp response.GetPermissionsResponse
	var dbRepo models.Repo
	ctrl.Tx.First(&dbRepo, repoId)
	if len(ctrl.Tx.GetErrors()) > 0 {
		resp.Success = false
		resp.Error = basicResponse.FatalError
		return ctrl.RenderJson(resp)
	}
	if dbRepo.ID != uint(repoId) {
		resp.Success = false
		resp.Error = basicResponse.NoRepositoryFoundError
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

	ctrl.Tx.First(&dbRepo, repoId)
	if len(ctrl.Tx.GetErrors()) > 0 {
		resp.Success = false
		resp.Error = basicResponse.FatalError
		revel.ERROR.Println(ctrl.Tx.GetErrors()[0].Error())
		return ctrl.RenderJson(resp)
	}
	if dbRepo.ID != uint(repoId) {
		resp.Success = false
		resp.Error = basicResponse.NoRepositoryFoundError
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
		ctrl.Tx.Create(&dbPermission)
	}

	if len(ctrl.Tx.GetErrors()) > 0 {
		resp.Success = false
		resp.Error = basicResponse.FatalError
		return ctrl.RenderJson(resp)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}
