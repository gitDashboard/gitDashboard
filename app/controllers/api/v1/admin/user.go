package admin

import (
	"github.com/gitDashboard/client/v1/admin/response"
	basicResponse "github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

type AdminUser struct {
	controllers.AdminController
}

func (ctrl *AdminUser) Search() revel.Result {
	var resp response.UsersResponse
	var query string
	ctrl.Params.Bind(&query, "query")
	revel.INFO.Println("Searching user query:", query)
	var dbUsers []models.User
	db := ctrl.Tx.Where("username LIKE ?", "%"+query+"%").Find(&dbUsers)
	if len(db.GetErrors()) > 0 {
		resp.Success = false
		resp.Error = basicResponse.FatalError
		for _, err := range db.GetErrors() {
			resp.Error.Message = resp.Error.Message + err.Error() + " "
		}
		return ctrl.RenderJson(resp)
	}
	resp.Users = make([]response.User, len(dbUsers), len(dbUsers))
	for i, dbUser := range dbUsers {
		resp.Users[i] = response.User{ID: dbUser.ID, Username: dbUser.Username, Type: dbUser.Type}
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminUser) List() revel.Result {
	var resp response.UsersResponse
	var dbUsers []models.User
	ctrl.Tx.Find(&dbUsers)
	resp.Users = make([]response.User, len(dbUsers), len(dbUsers))
	for i, dbUser := range dbUsers {
		resp.Users[i] = response.User{ID: dbUser.ID, Username: dbUser.Username, Type: dbUser.Type}
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}

func (ctrl *AdminUser) Save() revel.Result {
	var resp response.UserUpdateResponse
	var req response.User
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ctrl.RenderError(err)
	}
	var dbUser models.User
	if req.ID != 0 {
		//toUpdate
		ctrl.Tx.First(&dbUser, req.ID)
		if dbUser.ID != req.ID {
			resp.Success = false
			resp.Error = basicResponse.NoUserFoundError
			return ctrl.RenderJson(resp)
		}
	}
	dbUser.Username = req.Username
	dbUser.Type = req.Type
	if dbUser.Type == "internal" && req.Password != "" {
		dbUser.Password.Valid = true
		dbUser.Password.String = req.Password
	}
	var db *gorm.DB
	if req.ID != 0 {
		//update user
		db = ctrl.Tx.Save(&dbUser)
	} else {
		// new user
		db = ctrl.Tx.Create(&dbUser)
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

func (ctrl *AdminUser) Delete(userId uint) revel.Result {
	var resp response.UserDeleteResponse
	var dbUser models.User
	ctrl.Tx.First(&dbUser, userId)
	if dbUser.ID == 0 {
		resp.Success = false
		resp.Error = basicResponse.NoUserFoundError
		return ctrl.RenderJson(resp)
	}
	db := ctrl.Tx.Delete(&dbUser)
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