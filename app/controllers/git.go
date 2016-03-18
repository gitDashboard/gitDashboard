package controllers

import (
	"encoding/base64"
	"errors"
	"github.com/gitDashboard/client/v1/misc"
	"github.com/gitDashboard/gitDashboard/app/config"
	"github.com/revel/revel"
	"net/http"
	"net/http/cgi"
	"strings"
)

type GitRaw struct {
	*cgi.Handler
}

func (r GitRaw) Apply(req *revel.Request, resp *revel.Response) {
	r.ServeHTTP(resp.Out, req.Request)
}

type GitCtrl struct {
	GormController
}

func getCredentials(data string) (username, password string, err error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", "", err
	}
	strData := strings.Split(string(decodedData), ":")
	username = strData[0]
	password = strData[1]
	return
}

func (ctrl *GitCtrl) Repo() revel.Result {

	// The auth data is sent in the headers and will have a value of "Basic XXX" where XXX is base64 encoded data
	if auth := ctrl.Request.Header.Get("Authorization"); auth != "" {
		// Split up the string to get just the data, then get the credentials
		username, password, err := getCredentials(strings.Split(auth, " ")[1])
		if err != nil {
			return ctrl.RenderError(err)
		}
		dbUser, err := Login(ctrl.Tx, username, password)
		if dbUser == nil || err != nil {
			ctrl.Response.Status = http.StatusUnauthorized
			ctrl.Response.Out.Header().Set("WWW-Authenticate", `Basic realm="gitDashboard"`)
			return ctrl.RenderError(errors.New("401: Not authorized"))
		}
		pathInfo := strings.Replace(ctrl.Request.Request.URL.Path, "git/", "", 1)
		infoUrlPos := strings.LastIndex(pathInfo, "/info")
		var repoPath string
		if infoUrlPos != -1 {
			repoPath = pathInfo[1:infoUrlPos]
		} else {
			uploadPackUrlPos := strings.LastIndex(pathInfo, "/git-upload-pack")
			if uploadPackUrlPos != -1 {
				repoPath = pathInfo[1:uploadPackUrlPos]
			} else {
				receivePackUrlPos := strings.LastIndex(pathInfo, "/git-receive-pack")
				if receivePackUrlPos != -1 {
					repoPath = pathInfo[1:receivePackUrlPos]
				}
			}
		}
		authorized, locked, _ := CheckAutorization(ctrl.Tx, repoPath, username, "read", "/")
		eventId, err := StartEvent(ctrl.Tx, repoPath, "access", username, "/", pathInfo, misc.EventLevel_INFO, false)
		defer FinishEvent(ctrl.Tx, eventId)
		switch {
		case !authorized:
			StartEvent(ctrl.Tx, repoPath, "access", username, "/", "Forbidden", misc.EventLevel_WARN, true)
			ctrl.Response.Status = http.StatusForbidden
			return ctrl.RenderError(errors.New("403: Forbidden"))
		case locked:
			StartEvent(ctrl.Tx, repoPath, "access", username, "/", "Locked", misc.EventLevel_WARN, true)
			ctrl.Response.Status = http.StatusServiceUnavailable
			return ctrl.RenderError(errors.New("503 Service Unavailable"))
		}

		var resp GitRaw
		gitBackendCmd := cgi.Handler{
			Path:   config.GitHttpBackendPath(),
			Logger: revel.ERROR,
		}
		resp.Handler = &gitBackendCmd

		gitBackendCmd.Env = make([]string, 8)
		gitBackendCmd.Env[0] = "GIT_HTTP_EXPORT_ALL=true"
		gitBackendCmd.Env[1] = "GIT_PROJECT_ROOT=" + config.GitBasePath()
		gitBackendCmd.Env[2] = "PATH_INFO=" + pathInfo
		gitBackendCmd.Env[3] = "CONTENT_TYPE=" + ctrl.Request.Request.Header.Get("Content-Type")
		gitBackendCmd.Env[4] = "REQUEST_METHOD=" + ctrl.Request.Request.Method
		gitBackendCmd.Env[5] = "QUERY_STRING=" + ctrl.Request.Request.URL.Query().Encode()
		gitBackendCmd.Env[6] = "REMOTE_USER=" + username
		gitBackendCmd.Env[7] = "REMOTE_ADDR=" + ctrl.Request.RemoteAddr
		return resp
	} else {
		ctrl.Response.Status = http.StatusUnauthorized
		ctrl.Response.Out.Header().Set("WWW-Authenticate", `Basic realm="gitDashboard"`)
		return ctrl.RenderError(errors.New("401: Not authorized"))
	}
}
