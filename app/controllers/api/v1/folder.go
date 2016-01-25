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
	folderInfo.Admins = make([]response.ShortUserInfo, len(dbFolder.Admins), len(dbFolder.Admins))
	for a, admin := range dbFolder.Admins {
		folderInfo.Admins[a] = response.ShortUserInfo{ID: admin.ID, Username: admin.Username, Name: admin.Name}
		if admin.Email.Valid {
			folderInfo.Admins[a].Email = admin.Email.String
		}
	}
	folderInfo.ExtendedAdmins = folderInfo.Admins
	return folderInfo
}

func (ctrl *FolderCtrl) loadFolderAdmins(dbFolder *models.Folder, folderInfo *response.FolderInfo) error {
	if dbFolder.ParentID != 0 {
		var parentFolder models.Folder
		db := ctrl.Tx.Preload("Admins").Find(&parentFolder, dbFolder.ParentID)
		if db.Error != nil {
			return db.Error
		}
		for _, admin := range parentFolder.Admins {
			extAdminInfo := response.ShortUserInfo{ID: admin.ID, Username: admin.Username, Name: admin.Name}
			folderInfo.ExtendedAdmins = append(folderInfo.ExtendedAdmins, extAdminInfo)
		}
		return ctrl.loadFolderAdmins(&parentFolder, folderInfo)
	}
	return nil
}

func (ctrl *FolderCtrl) Get(folderId uint) revel.Result {
	var resp response.FolderGetResponse
	var dbFolder models.Folder
	db := ctrl.Tx.Preload("Admins").First(&dbFolder, folderId)
	if db.Error != nil {
		controllers.ErrorResp(&resp, response.DbError, db.Error)
		return ctrl.RenderJson(resp)
	}
	resp.Folder = getFolder(dbFolder)
	err := ctrl.loadFolderAdmins(&dbFolder, &resp.Folder)
	if err != nil {
		controllers.ErrorResp(&resp, response.DbError, err)
		return ctrl.RenderJson(resp)
	}
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
