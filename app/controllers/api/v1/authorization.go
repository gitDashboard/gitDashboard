package v1

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/gitDashboard/client/v1/misc"
	"github.com/gitDashboard/client/v1/request"
	"github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/auth"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type AuthorizationCtrl struct {
	controllers.GormController
}

func CheckAutorization(db *gorm.DB, repoDir, username, operation, refName string) (bool, error) {
	var user models.User
	//finding user
	db.Preload("Groups").Where("username = ?", username).First(&user)
	//searching for admin group
	isAdmin := false
	groupIds := make([]uint, len(user.Groups), len(user.Groups))
	for i, grp := range user.Groups {
		groupIds[i] = grp.ID
		if grp.Name == "admin" {
			isAdmin = true
		}
	}
	if isAdmin {
		return true, nil
	}
	//finding repo
	var repo models.Repo
	db.Where("path = ?", repoDir).First(&repo)
	//finding permissions
	var perms []models.Permission
	if len(groupIds) > 0 {
		db.Order("position").Where("repo_id = ? and (user_id = ? or group_id IN (?)) and type like ?", repo.ID, user.ID, groupIds, "%"+operation+"%").Find(&perms)
	} else {
		db.Order("position").Where("repo_id = ? and user_id = ? and type like ?", repo.ID, user.ID, "%"+operation+"%").Find(&perms)
	}
	authorized := false
	if len(perms) > 0 {
		for _, perm := range perms {
			if operation != "read" {
				match, err := regexp.MatchString(perm.Branch, refName)
				if err != nil {
					return false, err
				}
				if match {
					authorized = perm.Granted
				}
			} else {
				authorized = true
			}
		}
	}
	revel.INFO.Println("authorized:", authorized)
	return authorized, nil
}

func (c AuthorizationCtrl) CheckAuthorization() revel.Result {
	authReq := new(request.AuthorizationRequest)
	err := c.GetJSONBody(authReq)
	if err != nil {
		c.RenderError(err)
	}
	revel.INFO.Printf("CheckAuthorization parameter:%+v\n", authReq)

	authorized, err := CheckAutorization(c.Tx, authReq.RepositoryPath, authReq.Username, authReq.Operation, authReq.RefName)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(&response.AuthorizationResponse{Authorized: authorized})
}

func checkUserPassword(dbUser *models.User, password, userType string) error {
	if userType == "internal" {
		return bcrypt.CompareHashAndPassword([]byte(dbUser.Password.String), []byte(password))
	} else {
		return auth.Login(dbUser.Username, password)
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
		loginResp.Error = response.NoUserFoundError
		return ctrl.RenderJson(loginResp)
	}
	pwdErr := checkUserPassword(&dbUser, loginReq.Password, loginReq.Type)
	if pwdErr != nil {
		revel.ERROR.Println(pwdErr.Error())
		loginResp.Success = false
		loginResp.Error = response.AuthenticationFailedError
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
			loginResp.Error = response.FatalError
			loginResp.Error.Message = loginResp.Error.Message + err.Error()
		} else {
			loginResp.Success = true
			loginResp.JWT = jwtStr
		}
	}
	return ctrl.RenderJson(loginResp)
}
