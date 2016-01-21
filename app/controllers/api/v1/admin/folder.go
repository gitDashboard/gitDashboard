package admin

import (
	"github.com/gitDashboard/client/v1/admin/request"
	"github.com/gitDashboard/client/v1/admin/response"
	basicResponse "github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
	"os"
	"strings"
)

type AdminFolder struct {
	controllers.AdminController
}

func (ctrl *AdminFolder) CreateFolder() revel.Result {
	var req request.CreateFolderRequest
	var resp response.CreateFolderResponse
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		return ctrl.RenderError(err)
	}

	//searching parent
	var parentFolder *models.Folder
	if req.ParentID > 0 {
		parentFolder = new(models.Folder)
		db := ctrl.Tx.First(parentFolder, req.ParentID)
		if db.Error != nil {
			controllers.ErrorResp(&resp, basicResponse.DbError, db.Error)
			return ctrl.RenderJson(resp)
		}
	}
	dbFolder := models.Folder{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Description: req.Description,
	}
	var fullPath string
	if parentFolder != nil {
		dbFolder.Path = parentFolder.Path + "/" + req.Name
		fullPath = controllers.CleanSlashes(controllers.GitBasePath() + "/" + parentFolder.Path + "/" + req.Name)
	} else {
		dbFolder.Path = req.Name
		fullPath = controllers.CleanSlashes(controllers.GitBasePath() + "/" + req.Name)
	}
	//check if already exists on db
	var existDbFolder models.Folder
	db := ctrl.Tx.Where("path = ?", dbFolder.Path).First(&existDbFolder)
	if db.Error == nil && existDbFolder.ID > 0 {
		controllers.ErrorResp(&resp, basicResponse.AlreadyExistError, nil)
		return ctrl.RenderJson(resp)
	}

	db = ctrl.Tx.Create(&dbFolder)
	if db.Error != nil {
		controllers.ErrorResp(&resp, basicResponse.DbError, db.Error)
		return ctrl.RenderJson(resp)
	}
	revel.INFO.Println("creating folder", fullPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := os.Mkdir(fullPath, 0770)
		if err != nil {
			ctrl.Tx.Rollback()
			controllers.ErrorResp(&resp, basicResponse.FatalError, err)
		} else {
			resp.Success = true
		}
	} else {
		controllers.ErrorResp(&resp, basicResponse.AlreadyExistError, nil)
	}
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminFolder) Permissions(folderId uint) revel.Result {
	var resp response.GetPermissionsResponse
	var dbFolder models.Folder
	ctrl.Tx.First(&dbFolder, folderId)
	if len(ctrl.Tx.GetErrors()) > 0 {
		controllers.ErrorResp(&resp, basicResponse.FatalError, nil)
		return ctrl.RenderJson(resp)
	}
	if dbFolder.ID != uint(folderId) {
		controllers.ErrorResp(&resp, basicResponse.NoFolderFoundError, nil)
		return ctrl.RenderJson(resp)
	}
	ctrl.Tx.Preload("Users").Order("position").Where("folder_id=?", dbFolder.ID).Find(&dbFolder.Permissions)

	resp.Permissions = make([]response.Permission, len(dbFolder.Permissions), len(dbFolder.Permissions))
	for i, perm := range dbFolder.Permissions {
		var repoPerm response.Permission
		repoPerm.Ref = perm.Branch
		repoPerm.Position = perm.Position
		repoPerm.Granted = perm.Granted
		repoPerm.Users = make([]response.User, len(perm.Users), len(perm.Users))
		for u, user := range perm.Users {
			repoPerm.Users[u].Username = user.Username
			repoPerm.Users[u].Type = user.Type
			repoPerm.Users[u].Name = user.Name
			if user.Email.Valid {
				repoPerm.Users[u].Email = user.Email.String
			}
		}
		if perm.Type != "" {
			repoPerm.Types = strings.Split(perm.Type, ",")
		}
		resp.Permissions[i] = repoPerm
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminFolder) UpdatePermissions(folderId uint) revel.Result {
	var req request.UpdatePermissionsRequest
	var resp response.UpdatePermissionsResponse
	var dbFolder models.Folder
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ctrl.RenderError(err)
	}

	revel.INFO.Printf("UpdatePermissions req:%+v\n", req)

	db := ctrl.Tx.First(&dbFolder, folderId)
	if len(db.GetErrors()) > 0 {
		controllers.ErrorResp(&resp, basicResponse.FatalError, db.GetErrors()[0])
		return ctrl.RenderJson(resp)
	}
	if dbFolder.ID != uint(folderId) {
		controllers.ErrorResp(&resp, basicResponse.NoFolderFoundError, nil)
		return ctrl.RenderJson(resp)
	}

	//remove all old permissions
	ctrl.Tx.Where("folder_id = ?", folderId).Delete(models.Permission{})
	//insert all new permissions
	for _, newPerm := range req.Permissions {
		dbPermission := models.Permission{
			Branch:   newPerm.Ref,
			Granted:  newPerm.Granted,
			Position: newPerm.Position,
		}
		dbPermission.FolderID.Scan(folderId)
		for _, permType := range newPerm.Types {
			if permType != "" {
				dbPermission.Type = dbPermission.Type + "," + permType
			}
		}
		dbPermission.Type = dbPermission.Type[1:]
		dbPermission.Users = make([]models.User, len(newPerm.Users), len(newPerm.Users))
		for u, user := range newPerm.Users {
			var dbUser models.User
			ctrl.Tx.First(&dbUser, user.ID)
			dbPermission.Users[u] = dbUser
		}
		db := ctrl.Tx.Create(&dbPermission)
		if len(db.GetErrors()) > 0 {
			controllers.ErrorResp(&resp, basicResponse.FatalError, db.GetErrors()[0])
			return ctrl.RenderJson(resp)
		}
	}

	resp.Success = true
	return ctrl.RenderJson(resp)
}
