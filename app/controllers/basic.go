package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"io/ioutil"
)

type BasicController struct {
	*revel.Controller
}

func (c *BasicController) GetJSONBody(out interface{}) error {
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
