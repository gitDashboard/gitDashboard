package controllers

import (
	"github.com/gitDashboard/client/v1/response"
	"github.com/revel/revel"
	"regexp"
)

var SlashRegexp *regexp.Regexp

func init() {
	SlashRegexp = regexp.MustCompile("/{2,}")
}

func ErrorResp(resp response.IBasicResponse, respError response.Error, err error) {
	resp.SetSuccess(false)
	if err != nil {
		respError.Message = respError.Message + err.Error()
	}
	resp.SetError(respError)
	revel.ERROR.Println(err.Error())
}
