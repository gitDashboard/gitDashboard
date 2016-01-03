package v1

import (
	"fmt"
	"github.com/gitDashboard/client/v1/request"
	"github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
	git "gopkg.in/libgit2/git2go.v22"
	"io/ioutil"
	"os"
	"strings"
)

type RepoCtrl struct {
	controllers.JWTController
}

func (ctrl *RepoCtrl) getRepo(fullPath string) models.Repo {
	var repo models.Repo
	ctrl.Tx.Where("path = ? ", fullPath).First(&repo)
	return repo
}

func readRepoDescription(repoPath string, repoInfo *response.RepoInfo) {
	//read description
	descCnt, err := ioutil.ReadFile(repoPath + "/description")
	if err == nil {
		repoInfo.Description = string(descCnt)
	} else {
		revel.ERROR.Println(err.Error())
	}
}

func (ctrl *RepoCtrl) checkIsRepo(basePath, repoPath string, repoInfo *response.RepoInfo) error {
	fullRepoPath := controllers.CleanSlashes(repoPath + "/" + repoInfo.Name)
	repoInfo.Path = strings.Replace(fullRepoPath, basePath, "", 1)
	repo := ctrl.getRepo(fullRepoPath)
	if repo.ID > 0 {
		repoInfo.ID = repo.ID
		repoInfo.IsRepo = true
		readRepoDescription(fullRepoPath, repoInfo)
		//checking permission
		authorized, err := CheckAutorization(ctrl.Tx, fullRepoPath, ctrl.User.Username, "read", "")
		if err != nil {
			return err
		}
		repoInfo.IsAuthorized = authorized
	} else {
		repoInfo.Description = "Folder"
		repoInfo.IsAuthorized = true
	}

	return nil
}

func (ctrl *RepoCtrl) List() revel.Result {
	var req request.RepoListRequest
	var resp response.RepoListResponse
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		return ctrl.RenderError(err)
	}
	baseDirPath := revel.Config.StringDefault("git.baseDir", "/")
	currDirPath := controllers.CleanSlashes(baseDirPath + "/" + req.SubPath)
	revel.INFO.Println("Reading repositories from path:", currDirPath)

	repo := ctrl.getRepo(currDirPath)
	if repo.ID != 0 {
		//wrong ws , possible security problem: return empty list
		return ctrl.RenderJson(resp)
	}

	currDir, err := os.Open(currDirPath)
	if err != nil {
		return ctrl.RenderError(err)
	}

	finfo, err := currDir.Readdir(-1)
	if err != nil {
		return ctrl.RenderError(err)
	}
	resp.Repositories = make([]response.RepoInfo, 0, len(finfo))
	for _, f := range finfo {
		revel.INFO.Println("check ", f.Name())
		if f.IsDir() {
			var repoInfo response.RepoInfo
			repoInfo.Name = f.Name()

			err = ctrl.checkIsRepo(baseDirPath, currDirPath, &repoInfo)
			if err != nil {
				revel.ERROR.Println(err.Error())
				return ctrl.RenderError(err)
			}
			resp.Repositories = append(resp.Repositories, repoInfo)
		}
	}
	revel.INFO.Printf("result %+v\n ", resp)
	return ctrl.RenderJson(resp)
}

func (ctrl *RepoCtrl) Commits(repoId int) revel.Result {
	var dbRepo models.Repo
	var req request.RepoCommitsRequest
	var resp response.RepoCommitsResponse

	err := ctrl.GetJSONBody(&req)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		revel.ERROR.Println(err.Error())
		return ctrl.RenderJson(resp)
	}

	ctrl.Tx.First(&dbRepo, repoId)
	if dbRepo.ID != uint(repoId) {
		resp.Success = false
		resp.Error = response.NoRepositoryFoundError
		return ctrl.RenderJson(resp)
	}
	//checking permission
	authorized, err := CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	if !authorized {
		resp.Success = false
		resp.Error = response.PermissionDeniedError
		return ctrl.RenderJson(resp)
	}
	repo, err := git.OpenRepository(dbRepo.Path)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	refName := "refs/heads/" + req.Branch
	walk, err := repo.Walk()
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	resp.Success = true
	resp.Commits = make([]response.RepoCommit, 0)
	walk.Sorting(git.SortTopological | git.SortTime | git.SortReverse)

	walk.HideGlob("tags/*")
	end := req.Start + req.Count

	revStartPar := fmt.Sprintf("%s~%d", refName, req.Start)
	revel.INFO.Printf("revStart:%s\n", revStartPar)
	revStart, err := repo.Revparse(revStartPar)
	var revEndPar string
	var revEnd *git.Revspec
	if err != nil {
		goto end
	} else {
		walk.Push(revStart.From().Id())
	}
	revEndPar = fmt.Sprintf("%s~%d", refName, end)
	revel.INFO.Printf("revEnd:%s\n", revEndPar)
	revEnd, err = repo.Revparse(revEndPar)
	if err != nil {
		revel.WARN.Println("Requested range before first commit error:", err.Error())
	} else {
		walk.Hide(revEnd.From().Id())
	}

	walk.Iterate(func(commit *git.Commit) bool {
		var repoCmt response.RepoCommit
		repoCmt.Message = commit.Message()
		repoCmt.Author = commit.Author().Name
		repoCmt.Email = commit.Author().Email
		repoCmt.Date = commit.Author().When.UnixNano() / 1000000
		resp.Commits = append(resp.Commits, repoCmt)
		commit.Free()
		return true
	})
end:
	walk.Free()
	repo.Free()
	revel.INFO.Printf("response: %+v\n", resp)
	return ctrl.RenderJson(resp)
}

func (ctrl *RepoCtrl) Info(repoId int) revel.Result {
	var dbRepo models.Repo
	var resp response.RepoInfoResponse

	ctrl.Tx.First(&dbRepo, repoId)
	if dbRepo.ID != uint(repoId) {
		resp.Success = false
		resp.Error = response.NoRepositoryFoundError
		return ctrl.RenderJson(resp)
	}
	//checking permission
	authorized, err := CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	if !authorized {
		resp.Success = false
		resp.Error = response.PermissionDeniedError
		return ctrl.RenderJson(resp)
	}
	resp.Info.Path = "/" + strings.Replace(dbRepo.Path, revel.Config.StringDefault("git.baseDir", "/"), "", 1)
	resp.Info.FolderPath = resp.Info.Path[0:strings.LastIndex(resp.Info.Path, "/")]
	resp.Info.Url = revel.Config.StringDefault("git.baseUrl", "/") + "/" + resp.Info.Path
	resp.Info.ID = dbRepo.ID
	readRepoDescription(dbRepo.Path, &resp.Info)
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *RepoCtrl) Files(repoId int) revel.Result {
	var dbRepo models.Repo
	var resp response.RepoFilesResponse
	var req request.RepoFilesRequest

	err := ctrl.GetJSONBody(&req)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		revel.ERROR.Println(err.Error())
		return ctrl.RenderJson(resp)
	}

	ctrl.Tx.First(&dbRepo, repoId)
	if dbRepo.ID != uint(repoId) {
		resp.Success = false
		resp.Error = response.NoRepositoryFoundError
		return ctrl.RenderJson(resp)
	}
	//checking permission
	authorized, err := CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	if !authorized {
		resp.Success = false
		resp.Error = response.PermissionDeniedError
		return ctrl.RenderJson(resp)
	}
	repo, err := git.OpenRepository(dbRepo.Path)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	//find last commit
	rev, err := repo.Revparse(req.RefName)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		repo.Free()
		return ctrl.RenderJson(resp)
	}
	lastCommit, err := repo.LookupCommit(rev.From().Id())
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		repo.Free()
		return ctrl.RenderJson(resp)
	}
	var tree *git.Tree
	if req.Parent == "" {
		tree, err = lastCommit.Tree()
	} else {
		oid, _ := git.NewOid(req.Parent)
		tree, err = repo.LookupTree(oid)
	}
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		lastCommit.Free()
		repo.Free()
		return ctrl.RenderJson(resp)
	}
	resp.Files = make([]response.RepoFile, tree.EntryCount(), tree.EntryCount())

	for i := uint64(0); i < tree.EntryCount(); i++ {
		var respFile response.RepoFile
		fileEntry := tree.EntryByIndex(i)
		respFile.Id = fileEntry.Id.String()
		respFile.Name = fileEntry.Name
		if fileEntry.Type == git.ObjectTree {
			respFile.IsDir = true
		}
		resp.Files[i] = respFile
		revel.INFO.Printf("%+v\n", fileEntry)
	}
	resp.Success = true
	tree.Free()
	lastCommit.Free()
	repo.Free()
	return ctrl.RenderJson(resp)
}
