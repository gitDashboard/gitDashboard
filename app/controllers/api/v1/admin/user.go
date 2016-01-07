package admin

import (
	"github.com/gitDashboard/client/v1/admin/response"
	basicResponse "github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
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
