package controllers

import (
	"errors"
	"github.com/revel/revel"
)

type AdminController struct {
	JWTController
}

func (ctrl *AdminController) CheckPermission() revel.Result {
	if !ctrl.User.Admin {
		return ctrl.RenderError(errors.New("401: Not authorized"))
	} else {
		return nil
	}
}
