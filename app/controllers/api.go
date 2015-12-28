package controllers

import (
	"encoding/json"
	"io/ioutil"
)

type ApiController struct {
	GormController
}

func (c *ApiController) GetJSONBody(out interface{}) error {
	byteBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteBody, out)
	if err != nil {
		return err
	}
	return nil
}
