package v1

import (
	"github.com/gitDashboard/client/v1/request"
	"github.com/gitDashboard/client/v1/respond"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
	"regexp"
)

type AuthorizationCtrl struct {
	controllers.ApiController
}

func (c AuthorizationCtrl) CheckAuthorization() revel.Result {
	authReq := new(request.AuthorizationRequest)
	err := c.GetJSONBody(authReq)
	if err != nil {
		c.RenderError(err)
	}
	revel.INFO.Printf("CheckAuthorization parameter:%+v\n", authReq)

	var user models.User
	//finding user
	c.Tx.Where("username = ?", authReq.Username).First(&user)
	//finding repo
	var repo models.Repo
	c.Tx.Where("path = ?", authReq.RepositoryPath).First(&repo)
	//finding permissions
	var perms []models.Permission
	c.Tx.Where("repo_id = ? and user_id = ? and type=?", repo.ID, user.ID, authReq.Operation).Find(&perms)
	authorized := false
	for _, perm := range perms {
		revel.INFO.Printf("perm: %v", perm)
		match, err := regexp.MatchString(perm.Branch, authReq.RefName)
		if err != nil {
			revel.ERROR.Println("Error checking permission regex:", perm.Branch, err.Error())
		}
		if match {
			authorized = perm.Granted
			if authorized {
				break
			}
		}
	}

	return c.RenderJson(&respond.AuthorizationRespond{Authorized: authorized})
}
