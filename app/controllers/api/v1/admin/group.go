package admin

import (
	"github.com/gitDashboard/client/v1/admin/response"
	basicResponse "github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

type AdminGroup struct {
	controllers.AdminController
}

func (ctrl *AdminGroup) Search() revel.Result {
	var resp response.GroupsResponse
	var query string
	ctrl.Params.Bind(&query, "query")
	revel.INFO.Println("Searching group query:", query)
	var dbGroups []models.Group
	db := ctrl.Tx.Preload("Users").Where("name LIKE ?", "%"+query+"%").Find(&dbGroups)
	if len(db.GetErrors()) > 0 {
		resp.Success = false
		resp.Error = basicResponse.FatalError
		for _, err := range db.GetErrors() {
			resp.Error.Message = resp.Error.Message + err.Error() + " "
		}
		return ctrl.RenderJson(resp)
	}
	resp.Groups = make([]response.Group, len(dbGroups), len(dbGroups))
	for i, dbGroup := range dbGroups {
		resp.Groups[i] = response.Group{ID: dbGroup.ID, Name: dbGroup.Name, Description: dbGroup.Description}
		resp.Groups[i].Users = make([]response.User, len(dbGroup.Users), len(dbGroup.Users))
		for u, dbUser := range dbGroup.Users {
			resp.Groups[i].Users[u] = response.User{ID: dbUser.ID, Username: dbUser.Username, Type: dbUser.Type}
		}
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminGroup) List() revel.Result {
	var resp response.GroupsResponse
	var dbGroups []models.Group
	ctrl.Tx.Preload("Users").Find(&dbGroups)
	resp.Groups = make([]response.Group, len(dbGroups), len(dbGroups))
	for i, dbGroup := range dbGroups {
		resp.Groups[i] = response.Group{ID: dbGroup.ID, Name: dbGroup.Name, Description: dbGroup.Description}
		resp.Groups[i].Users = make([]response.User, len(dbGroup.Users), len(dbGroup.Users))
		for u, dbUser := range dbGroup.Users {
			resp.Groups[i].Users[u] = response.User{ID: dbUser.ID, Username: dbUser.Username, Type: dbUser.Type}
		}
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminGroup) Save() revel.Result {
	var resp response.GroupUpdateResponse
	var req response.Group
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ctrl.RenderError(err)
	}
	var dbGroup models.Group
	if req.ID != 0 {
		//toUpdate
		ctrl.Tx.First(&dbGroup, req.ID)
		if dbGroup.ID != req.ID {
			resp.Success = false
			resp.Error = basicResponse.NoGroupFoundError
			return ctrl.RenderJson(resp)
		}
	}
	if req.Name != "admin" {
		dbGroup.Name = req.Name
		dbGroup.Description = req.Description
	}
	ctrl.Tx.Model(&dbGroup).Association("Users").Clear()
	for _, grpUser := range req.Users {
		ctrl.Tx.Model(&dbGroup).Association("Users").Append(models.User{ID: grpUser.ID})
	}

	var db *gorm.DB
	if req.ID != 0 {
		//update group
		var grp models.Group
		db = ctrl.Tx.Model(&grp).Where("id = ?", req.ID).Updates(map[string]interface{}{"name": dbGroup.Name, "description": dbGroup.Description})
	} else {
		// new user
		db = ctrl.Tx.Create(&dbGroup)
	}
	if len(db.GetErrors()) > 0 {
		resp.Success = false
		resp.Error = basicResponse.FatalError
		for _, err := range db.GetErrors() {
			resp.Error.Message = resp.Error.Message + err.Error() + " "
		}
		return ctrl.RenderJson(resp)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminGroup) Delete(userId uint) revel.Result {
	var resp response.GroupDeleteResponse
	var dbGroup models.Group
	ctrl.Tx.First(&dbGroup, userId)
	if dbGroup.ID == 0 {
		resp.Success = false
		resp.Error = basicResponse.NoGroupFoundError
		return ctrl.RenderJson(resp)
	}
	db := ctrl.Tx.Delete(&dbGroup)
	if len(db.GetErrors()) > 0 {
		resp.Success = false
		resp.Error = basicResponse.FatalError
		for _, err := range db.GetErrors() {
			resp.Error.Message = resp.Error.Message + err.Error() + " "
		}
		return ctrl.RenderJson(resp)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}
