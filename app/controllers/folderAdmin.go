package controllers

import (
	"errors"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
)

type FolderAdminController struct {
	JWTController
	Folder models.Folder
}

func (ctrl *FolderAdminController) isAuthorized(folderToCheck *models.Folder) bool {
	if folderToCheck.ID != 0 {
		for _, admin := range folderToCheck.Admins {
			if admin.ID == ctrl.User.ID {
				return true
			}
		}
		var parentFolder models.Folder
		db := ctrl.Tx.Preload("Admins").Find(&parentFolder, folderToCheck.ParentID)
		if db.Error != nil {
			revel.ERROR.Println(db.Error)
			return false
		}
		return ctrl.isAuthorized(&parentFolder)
	}
	return false
}

func (ctrl *FolderAdminController) CheckPermission() revel.Result {
	if ctrl.User.Admin {
		return nil
	} else {
		var folderId uint
		var repoId uint
		authorized := false
		ctrl.Params.Bind(&folderId, "folderId")
		ctrl.Params.Bind(&repoId, "repoId")
		if folderId == 0 && repoId != 0 {
			var repo models.Repo
			db := ctrl.Tx.First(&repo, repoId)
			if db.Error != nil {
				revel.ERROR.Println(db.Error)
				return ctrl.RenderError(db.Error)
			}
			folderId = repo.FolderID
		}
		if folderId != 0 {
			db := ctrl.Tx.Preload("Admins").First(&ctrl.Folder, folderId)
			if db.Error != nil {
				revel.ERROR.Println(db.Error)
				return ctrl.RenderError(db.Error)
			}
			authorized = ctrl.isAuthorized(&ctrl.Folder)

		}
		if !authorized {
			return ctrl.RenderError(errors.New("401: Not authorized"))
		}
		return nil
	}
}
