package controllers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gitDashboard/client/v1/misc"
	"github.com/mitchellh/mapstructure"
	"github.com/revel/revel"
	"net/http"
)

type JWTController struct {
	GormController
	User misc.JWTUser
}

func (ctrl *JWTController) ParseToken() revel.Result {
	token, err := jwt.ParseFromRequest(ctrl.Request.Request, func(t *jwt.Token) (interface{}, error) {
		return []byte(revel.Config.StringDefault("jwt.secret", "")), nil
	})

	if err != nil {
		revel.ERROR.Println(err)
		ctrl.Response.Status = http.StatusUnauthorized
		return ctrl.RenderError(errors.New("401: Not authorized"))
	} else {
		err := mapstructure.Decode(token.Claims, &ctrl.User)
		if err != nil {
			revel.ERROR.Println(err)
		}
		return nil
	}
}
