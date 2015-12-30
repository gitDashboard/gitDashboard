package v1

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/gitDashboard/client/v1/misc"
	"github.com/gitDashboard/client/v1/request"
	"github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type AuthorizationCtrl struct {
	controllers.GormController
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

	return c.RenderJson(&response.AuthorizationResponse{Authorized: authorized})
}

func checkUserPassword(dbUser *models.User, password, userType string) error {
	if userType == "internal" {
		return bcrypt.CompareHashAndPassword([]byte(dbUser.Password.String), []byte(password))
	} else {
		//TODO : LDAP
		return nil
	}
}

func (ctrl *AuthorizationCtrl) Login() revel.Result {
	var loginReq request.LoginRequest
	var loginResp response.LoginResponse
	ctrl.GetJSONBody(&loginReq)
	revel.INFO.Printf("Login req:%+v\n", loginReq)
	var dbUser models.User
	ctrl.Tx.Preload("Groups").Where("username = ? and type = ?", loginReq.Username, loginReq.Type).First(&dbUser)

	if dbUser.ID == 0 {
		//no user found
		loginResp.Success = false
		loginResp.Message = "No user found with username:" + loginReq.Username
	}
	pwdErr := checkUserPassword(&dbUser, loginReq.Password, loginReq.Type)
	if pwdErr != nil {
		loginResp.Success = false
		loginResp.Message = "Login failed for user:" + loginReq.Username
	} else {
		var jwtUser misc.JWTUser
		jwtUser.Username = loginReq.Username
		jwtUser.Groups = make([]string, len(dbUser.Groups), len(dbUser.Groups))
		for i, group := range dbUser.Groups {
			jwtUser.Groups[i] = group.Name
		}

		jwtToken := jwt.New(jwt.SigningMethodHS256)
		jwtToken.Claims = structs.Map(jwtUser)
		jwtToken.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		jwtStr, err := jwtToken.SignedString([]byte(revel.Config.StringDefault("jwt.secret", "")))
		if err != nil {
			revel.ERROR.Println(err.Error())
			revel.ERROR.Printf("%+v\n", jwtToken)
			loginResp.Success = false
			loginResp.Message = err.Error()
		} else {
			loginResp.Success = true
			loginResp.JWT = jwtStr
		}
	}
	return ctrl.RenderJson(loginResp)
}
