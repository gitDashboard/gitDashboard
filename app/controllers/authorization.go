package controllers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/gitDashboard/client/v1/misc"
	"github.com/gitDashboard/client/v1/request"
	"github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/auth"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/gitDashboard/gitDashboard/app/repoManager"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type AuthorizationCtrl struct {
	GormController
}

func checkFolderAuthorization(db *gorm.DB, folderId, userID uint, operation, refName string) bool {
	var dbFolder models.Folder
	dbExec := db.First(&dbFolder, folderId)
	if dbExec.Error != nil {
		revel.ERROR.Println(dbExec.Error)
		return false
	}
	var authorized bool
	var perms []models.Permission
	db.Joins("inner join users_permissions on users_permissions.permission_id = permissions.id").Order("position").Where("folder_id = ? and users_permissions.user_id = ? and type like ?", folderId, userID, "%"+operation+"%").Find(&perms)

	oneMatch := false
	for _, perm := range perms {
		if operation != "read" {
			revel.INFO.Printf("checking permission : %+v\n", perm)
			match, err := regexp.MatchString(perm.Branch, refName)
			if err != nil {
				return false
			}
			if match {
				oneMatch = true
				revel.INFO.Println("match found")
				authorized = perm.Granted
			} else {
				revel.INFO.Println("match not found")
			}
		} else {
			oneMatch = true
			authorized = true
		}
	}
	if !oneMatch && dbFolder.ParentID != 0 {
		revel.INFO.Println("checkAutentication : no match on folder, search on parent folder")
		authorized = checkFolderAuthorization(db, dbFolder.ParentID, userID, operation, refName)
	}
	return authorized
}

func CheckAutorization(db *gorm.DB, repoDir, username, operation, refName string) (bool, bool, error) {
	var user models.User
	//finding user
	db.Where("username = ?", username).First(&user)

	//finding repo
	var repo models.Repo
	db.Where("path = ?", repoDir).First(&repo)

	if user.Admin {
		return true, repo.Locked, nil
	}
	//finding permissions
	var perms []models.Permission
	db.Joins("inner join users_permissions on users_permissions.permission_id = permissions.id").Order("position").Where("repo_id = ? and users_permissions.user_id = ? and type like ?", repo.ID, user.ID, "%"+operation+"%").Find(&perms)
	authorized := false
	oneMatch := false
	if len(perms) > 0 {
		for _, perm := range perms {
			if operation != "read" {
				match, err := regexp.MatchString(perm.Branch, refName)
				if err != nil {
					return false, false, err
				}
				if match {
					oneMatch = true
					authorized = perm.Granted
				}
			} else {
				oneMatch = true
				authorized = true
			}
		}
	}
	if !oneMatch {
		//search on parent folder authorization
		revel.INFO.Println("checkAutentication : no match on repository, search on parent folder")
		authorized = checkFolderAuthorization(db, repo.FolderID, user.ID, operation, refName)
	}

	revel.INFO.Println("authorized:", authorized)
	return authorized, repo.Locked, nil
}

func (c AuthorizationCtrl) CheckAuthorization() revel.Result {
	authReq := new(request.AuthorizationRequest)
	err := c.GetJSONBody(authReq)
	if err != nil {
		c.RenderError(err)
	}
	revel.INFO.Printf("CheckAuthorization parameter:%+v\n", authReq)

	authorized, locked, err := CheckAutorization(c.Tx, authReq.RepositoryPath, authReq.Username, authReq.Operation, authReq.RefName)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(&response.AuthorizationResponse{Authorized: authorized, Locked: locked})
}

func checkUserPassword(dbUser *models.User, password, userType string) error {
	if userType == "internal" {
		return bcrypt.CompareHashAndPassword([]byte(dbUser.Password.String), []byte(password))
	} else {
		return auth.Login(dbUser.Username, password)
	}
}

func Login(db *gorm.DB, username, password string) (*models.User, error) {
	var dbUser models.User
	db.Where("username = ?", username).First(&dbUser)
	if dbUser.ID == 0 {
		return nil, nil
	}
	return &dbUser, checkUserPassword(&dbUser, password, dbUser.Type)
}

func (ctrl *AuthorizationCtrl) Login() revel.Result {
	var loginReq request.LoginRequest
	var loginResp response.LoginResponse
	ctrl.GetJSONBody(&loginReq)
	revel.INFO.Printf("Login req:%+v\n", loginReq)
	dbUser, loginErr := Login(ctrl.Tx, loginReq.Username, loginReq.Password)
	if dbUser == nil {
		//no user found
		loginResp.Success = false
		loginResp.Error = response.NoUserFoundError
		return ctrl.RenderJson(loginResp)
	}
	if loginErr != nil {
		revel.ERROR.Println(loginErr.Error())
		loginResp.Success = false
		loginResp.Error = response.AuthenticationFailedError
	} else {
		var jwtUser misc.JWTUser
		jwtUser.ID = dbUser.ID
		jwtUser.Username = loginReq.Username
		jwtUser.Name = dbUser.Name
		jwtUser.Email = dbUser.Email.String
		jwtUser.Admin = dbUser.Admin

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

func StartEvent(db *gorm.DB, repoPath, eventType, username, reference, description string, level misc.EventLevel, finished bool) (uint, error) {
	repo, err := repoManager.GetRepo(db, repoPath)
	if err != nil {
		return 0, err
	}
	if repo.ID <= 0 {
		return 0, errors.New("No Repository found")
	}
	var event *models.Event
	if !finished {
		event, err = repoManager.StartRepoEvent(db, repo.ID, eventType, username, reference, description, level)
	} else {
		event, err = repoManager.AddRepoEvent(db, repo.ID, eventType, username, reference, description, level)
	}
	return event.ID, err
}

func (ctrl *AuthorizationCtrl) StartEvent(finished bool) revel.Result {
	var resp response.RepoEventResponse
	var req request.RepoEventRequest
	err := ctrl.GetJSONBody(&req)
	if err != nil {
		ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	eventID, err := StartEvent(ctrl.Tx, req.RepositoryPath, req.Type, req.User, req.Reference, req.Description, req.Level, finished)
	if err != nil {
		ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}

	resp.Success = true
	resp.EventID = eventID
	return ctrl.RenderJson(resp)
}

func FinishEvent(db *gorm.DB, eventId uint) error {
	return repoManager.FinishRepoEvent(db, eventId)
}

func (ctrl *AuthorizationCtrl) FinishEvent(eventId uint) revel.Result {
	var resp response.BasicResponse

	err := FinishEvent(ctrl.Tx, eventId)
	if err != nil {
		ErrorResp(&resp, response.FatalError, err)
		return ctrl.RenderJson(resp)
	}
	resp.Success = true
	return ctrl.RenderJson(resp)
}
