package v1

import (
	"github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
)

type FolderCtrl struct {
	controllers.JWTController
}

func getFolder(dbFolder models.Folder) response.FolderInfo {
	var folderInfo response.FolderInfo
	folderInfo.ID = dbFolder.ID
	folderInfo.Name = dbFolder.Name
	folderInfo.ParentID = dbFolder.ParentID
	folderInfo.Description = dbFolder.Description
	folderInfo.Path = dbFolder.Path
	return folderInfo
}

func (ctrl *FolderCtrl) Get(folderId uint) revel.Result {
	var resp response.FolderGetResponse
	var dbFolder models.Folder
	db := ctrl.Tx.First(&dbFolder, folderId)
	if db.Error != nil {
		controllers.ErrorResp(&resp, response.DbError, db.Error)
		return ctrl.RenderJson(resp)
	}
	resp.Folder = getFolder(dbFolder)
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *FolderCtrl) List(parentFolder uint) revel.Result {
	var resp response.FolderListResponse
	var folders []models.Folder
	db := ctrl.Tx.Where("parent_id = ?", parentFolder).Order("name asc").Find(&folders)

	if db.Error != nil {
		controllers.ErrorResp(&resp, response.DbError, db.Error)
		return ctrl.RenderJson(resp)
	}

	resp.Folders = make([]response.FolderInfo, len(folders), len(folders))
	for i, folder := range folders {
		resp.Folders[i] = getFolder(folder)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}
