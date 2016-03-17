package v1

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gitDashboard/client"
	"github.com/gitDashboard/client/v1/request"
	"github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
	git "gopkg.in/libgit2/git2go.v22"
	"html/template"
	"io/ioutil"
	"path"
	"regexp"
	"sort"
	"strings"
)

type RepoCtrl struct {
	controllers.JWTController
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

func (ctrl *RepoCtrl) List(parentFolder uint) revel.Result {
	var resp response.RepoListResponse
	var repos []models.Repo
	ctrl.Tx.Where("folder_id = ?", parentFolder).Order("name asc").Find(&repos)

	resp.Repositories = make([]response.RepoInfo, len(repos), len(repos))
	for i, repo := range repos {
		var repoInfo response.RepoInfo
		repoInfo.ID = repo.ID
		repoInfo.Name = repo.Name
		readRepoDescription(controllers.CleanSlashes(controllers.GitBasePath()+"/"+repo.Path), &repoInfo)
		repoInfo.Path = repo.Path
		repoInfo.Url = revel.Config.StringDefault("git.baseUrl", "/") + controllers.CleanSlashes(repoInfo.Path)
		repoInfo.IsAuthorized, repoInfo.Locked, _ = controllers.CheckAutorization(ctrl.Tx, repo.Path, ctrl.User.Username, "read", "/")
		resp.Repositories[i] = repoInfo
	}
	return ctrl.RenderJson(resp)
}

func (ctrl *RepoCtrl) Commits(repoId int) revel.Result {
	var dbRepo models.Repo
	var req request.RepoCommitsRequest
	var resp response.RepoCommitsResponse

	err := ctrl.GetJSONBody(&req)
	if err != nil {
		controllers.ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}

	ctrl.Tx.First(&dbRepo, repoId)
	if dbRepo.ID != uint(repoId) {
		resp.Success = false
		resp.Error = response.NoRepositoryFoundError
		return ctrl.RenderJson(resp)
	}
	//checking permission
	authorized, _, err := controllers.CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
	if err != nil {
		controllers.ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	if !authorized {
		controllers.ErrorResp(&resp, response.PermissionDeniedError, nil)
		return ctrl.RenderJson(resp)
	}
	repo, err := git.OpenRepository(controllers.CleanSlashes(controllers.GitBasePath() + "/" + dbRepo.Path))
	if err != nil {
		controllers.ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	defer repo.Free()
	refName := req.Branch
	walk, err := repo.Walk()
	if err != nil {
		controllers.ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	defer walk.Free()
	resp.Success = true
	resp.Commits = make([]response.RepoCommit, 0)
	walk.Sorting(git.SortTopological | git.SortTime)

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
		repoCmt.ID = commit.Id().String()
		repoCmt.Message = commit.Message()
		repoCmt.Author = commit.Author().Name
		repoCmt.Email = commit.Author().Email
		repoCmt.Date = commit.Author().When.UnixNano() / 1000000
		resp.Commits = append(resp.Commits, repoCmt)
		commit.Free()
		return true
	})
end:
	return ctrl.RenderJson(resp)
}

func (ctrl *RepoCtrl) Info(repoId int) revel.Result {
	var dbRepo models.Repo
	var resp response.RepoInfoResponse

	ctrl.Tx.First(&dbRepo, repoId)
	if dbRepo.ID != uint(repoId) {
		controllers.ErrorResp(&resp, response.NoRepositoryFoundError, nil)
		return ctrl.RenderJson(resp)
	}
	//checking permission
	authorized, _, err := controllers.CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
	if err != nil {
		controllers.ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	if !authorized {
		controllers.ErrorResp(&resp, response.PermissionDeniedError, nil)
		return ctrl.RenderJson(resp)
	}

	repo, err := git.OpenRepository(controllers.CleanSlashes(controllers.GitBasePath() + "/" + dbRepo.Path))
	if err != nil {
		controllers.ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	defer repo.Free()
	refIt, err := repo.NewReferenceIterator()
	if err != nil {
		controllers.ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	defer refIt.Free()
	refNameIt := refIt.Names()
	refName, refNameErr := refNameIt.Next()
	for refNameErr == nil {
		resp.Info.References = append(resp.Info.References, refName)
		refName, refNameErr = refNameIt.Next()
	}
	resp.Info.FolderID = dbRepo.FolderID
	resp.Info.Name = path.Base(dbRepo.Path)
	resp.Info.Path = dbRepo.Path
	resp.Info.Url = revel.Config.StringDefault("git.baseUrl", "/") + controllers.CleanSlashes(resp.Info.Path)
	resp.Info.ID = dbRepo.ID
	resp.Info.Locked = dbRepo.Locked
	readRepoDescription(controllers.CleanSlashes(controllers.GitBasePath()+"/"+dbRepo.Path), &resp.Info)
	resp.Success = true
	return ctrl.RenderJson(resp)
}

type ByFolder []response.RepoFile

func (a ByFolder) Len() int      { return len(a) }
func (a ByFolder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByFolder) Less(i, j int) bool {
	var fakeStr1 string
	var fakeStr2 string
	if a[i].IsDir {
		fakeStr1 = "d" + a[i].Name
	} else {
		fakeStr1 = "f" + a[i].Name
	}
	if a[j].IsDir {
		fakeStr2 = "d" + a[j].Name
	} else {
		fakeStr2 = "f" + a[j].Name
	}
	return fakeStr1 < fakeStr2
}

func (ctrl *RepoCtrl) Files(repoId int) revel.Result {
	var dbRepo models.Repo
	var resp response.RepoFilesResponse
	var req request.RepoFilesRequest

	err := ctrl.GetJSONBody(&req)
	if err != nil {
		controllers.ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}

	ctrl.Tx.First(&dbRepo, repoId)
	if dbRepo.ID != uint(repoId) {
		resp.Success = false
		resp.Error = response.NoRepositoryFoundError
		return ctrl.RenderJson(resp)
	}
	//checking permission
	authorized, _, err := controllers.CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
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
	repo, err := git.OpenRepository(controllers.CleanSlashes(controllers.GitBasePath() + "/" + dbRepo.Path))
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	defer repo.Free()
	//find last commit
	rev, err := repo.Revparse(req.RefName)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	lastCommit, err := repo.LookupCommit(rev.From().Id())
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	defer lastCommit.Free()
	var tree *git.Tree
	if req.Parent == "" {
		tree, err = lastCommit.Tree()
	} else {
		oid, _ := git.NewOid(req.Parent)
		tree, err = repo.LookupTree(oid)
		if err != nil {
			resp.ParentTreeId = req.Parent
		}
	}
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	defer tree.Free()
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
	sort.Sort(ByFolder(resp.Files))
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *RepoCtrl) FileContent(repoId int, fileRef string) revel.Result {
	var dbRepo models.Repo
	var resp response.RepoFileContentResponse
	revel.INFO.Printf("FileContent repoId:%d fileRef:%s\n", repoId, fileRef)

	ctrl.Tx.First(&dbRepo, repoId)
	if dbRepo.ID != uint(repoId) {
		resp.Success = false
		resp.Error = response.NoRepositoryFoundError
		return ctrl.RenderJson(resp)
	}
	//checking permission
	authorized, _, err := controllers.CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
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
	repo, err := git.OpenRepository(controllers.CleanSlashes(controllers.GitBasePath() + "/" + dbRepo.Path))
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	defer repo.Free()
	fileOid, err := git.NewOid(fileRef)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	blobFile, err := repo.LookupBlob(fileOid)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	defer blobFile.Free()

	resp.Success = true
	resp.Size = blobFile.Size()
	resp.Content = base64.StdEncoding.EncodeToString(blobFile.Contents())
	return ctrl.RenderJson(resp)
}

func Diff2Html(diffContent string) string {
	var result string
	rowRegexp := regexp.MustCompile("^@@.*?@@")
	lines := strings.Split(diffContent, "\n")
	for _, line := range lines {

		htmlLine := template.HTMLEscapeString(line)
		htmlLine = strings.Replace(htmlLine, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
		switch {
		case strings.HasPrefix(line, "+++"):
			result += "<p class=\"oldFile\">" + htmlLine + "</p>"
		case strings.HasPrefix(line, "---"):
			result += "<p class=\"newFile\">" + htmlLine + "</p>"
		case strings.HasPrefix(line, "+"):
			result += "<p class=\"inserted\">" + htmlLine + "</p>"
		case strings.HasPrefix(line, "-"):
			result += "<p class=\"removed\">" + htmlLine + "</p>"
		case strings.HasPrefix(line, "diff"):
			result += "<p class=\"diffCommand\">" + htmlLine + "</p>"
		default:
			if rowRegexp.MatchString(htmlLine) {
				matches := rowRegexp.FindAllString(htmlLine, -1)
				htmlLine = "<span class=\"linenumber\">" + matches[0] + "</span>" + htmlLine[len(matches[0]):]
			}
			result += "<p>" + htmlLine + "</p>"
		}
	}
	return result
}

func (ctrl *RepoCtrl) Commit(repoId uint, commitId string) revel.Result {
	var dbRepo models.Repo
	var resp response.RepoCommitResponse
	revel.INFO.Printf("Commit request repoId:%d commitId:%s\n", repoId, commitId)
	db := ctrl.Tx.First(&dbRepo, repoId)
	if len(db.GetErrors()) > 0 {
		resp.Success = false
		resp.Error = response.FatalError
		for _, err := range db.GetErrors() {
			resp.Error.Message = resp.Error.Message + err.Error()
		}
	}
	//checking permission
	authorized, _, err := controllers.CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
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
	repo, err := git.OpenRepository(controllers.CleanSlashes(controllers.GitBasePath() + "/" + dbRepo.Path))
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	defer repo.Free()
	commitOid, _ := git.NewOid(commitId)
	commit, err := repo.LookupCommit(commitOid)
	if err != nil {
		resp.Success = false
		resp.Error = response.FatalError
		resp.Error.Message = resp.Error.Message + err.Error()
		return ctrl.RenderJson(resp)
	}
	defer commit.Free()
	resp.Success = true
	resp.Commit.ID = commitId
	resp.Commit.Message = commit.Message()
	resp.Commit.Author = commit.Author().Name
	resp.Commit.Email = commit.Author().Email
	resp.Commit.Date = commit.Author().When.UnixNano() / 1000000

	if commit.ParentCount() == 1 {
		parentCommit := commit.Parent(0)
		defer parentCommit.Free()
		diffOpts, _ := git.DefaultDiffOptions()
		parentTree, _ := parentCommit.Tree()
		currTree, _ := commit.Tree()
		diff, err := repo.DiffTreeToTree(parentTree, currTree, &diffOpts)
		if err != nil {
			resp.Success = false
			resp.Error = response.FatalError
			resp.Error.Message = resp.Error.Message + err.Error()
			return ctrl.RenderJson(resp)
		}
		defer diff.Free()
		nDelta, err := diff.NumDeltas()
		if err != nil {
			resp.Success = false
			resp.Error = response.FatalError
			resp.Error.Message = resp.Error.Message + err.Error()
			return ctrl.RenderJson(resp)
		}
		for d := 0; d < nDelta; d++ {
			delta, err := diff.GetDelta(d)
			if err != nil {
				resp.Success = false
				resp.Error = response.FatalError
				resp.Error.Message = resp.Error.Message + err.Error()
				return ctrl.RenderJson(resp)
			}
			patch, err := diff.Patch(d)
			if err != nil {
				resp.Success = false
				resp.Error = response.FatalError
				resp.Error.Message = resp.Error.Message + err.Error()
				return ctrl.RenderJson(resp)
			}
			defer patch.Free()

			var diffFile response.RepoDiffFile
			patchContent, err := patch.String()

			if err != nil {
				resp.Success = false
				resp.Error = response.FatalError
				resp.Error.Message = resp.Error.Message + err.Error()
				return ctrl.RenderJson(resp)
			}
			diffFile.Patch = Diff2Html(patchContent)

			if delta.Status == git.DeltaAdded {
				diffFile.Type = "added"
			}
			if delta.Status == git.DeltaDeleted {
				diffFile.Type = "deleted"
			}
			if delta.Status == git.DeltaModified {
				diffFile.Type = "modified"
			}
			if delta.Status == git.DeltaRenamed {
				diffFile.Type = "renamed"
			}
			diffFile.OldId = delta.OldFile.Oid.String()
			diffFile.OldName = delta.OldFile.Path
			diffFile.NewId = delta.NewFile.Oid.String()
			diffFile.NewName = delta.NewFile.Path
			resp.Files = append(resp.Files, diffFile)
		}
	}
	return ctrl.RenderJson(resp)
}

func (ctrl *RepoCtrl) Graph(repoId uint) revel.Result {
	var dbRepo models.Repo
	revel.INFO.Printf("Graph repoId:%d\n", repoId)
	query := ctrl.Tx.First(&dbRepo, repoId)
	if query.Error != nil {
		return ctrl.RenderError(query.Error)
	}
	//checking permission
	authorized, _, err := controllers.CheckAutorization(ctrl.Tx, dbRepo.Path, ctrl.User.Username, "read", "")
	if err != nil {
		return ctrl.RenderError(err)
	}
	if !authorized {
		return ctrl.RenderError(errors.New("401: Not authorized"))
	}
	jsGraph, err := client.GenerateGraph(controllers.CleanSlashes(controllers.GitBasePath() + "/" + dbRepo.Path))
	if err != nil {
		return ctrl.RenderError(err)
	}
	revel.INFO.Println(jsGraph)
	return ctrl.RenderText(jsGraph)
}
