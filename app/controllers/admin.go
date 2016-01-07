package controllers

import (
	"errors"
	"github.com/revel/revel"
)

type AdminController struct {
	JWTController
}

func (ctrl *AdminController) CheckPermission() revel.Result {
	isAdmin := false
	for _, grp := range ctrl.User.Groups {
		if grp == "admin" {
			isAdmin = true
		}
	}

	if !isAdmin {
		return ctrl.RenderError(errors.New("401: Not authorized"))
	} else {
		return nil
	}
}
