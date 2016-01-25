package v1

import (
	"github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
)

type UserCtrl struct {
	controllers.JWTController
}

func (ctrl *UserCtrl) Search() revel.Result {
	var resp response.UsersResponse
	var query string
	ctrl.Params.Bind(&query, "query")
	revel.INFO.Println("Searching user query:", query)
	var dbUsers []models.User
	db := ctrl.Tx.Where("username LIKE ?", "%"+query+"%").Find(&dbUsers)
	if db.Error != nil {
		controllers.ErrorResp(&resp, response.DbError, db.Error)
		return ctrl.RenderJson(resp)
	}
	resp.Users = make([]response.User, len(dbUsers), len(dbUsers))
	for i, dbUser := range dbUsers {
		resp.Users[i] = response.User{ID: dbUser.ID, Username: dbUser.Username, Name: dbUser.Name, Type: dbUser.Type, Admin: dbUser.Admin}
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}
