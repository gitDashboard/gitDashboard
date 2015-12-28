package models

import (
	"github.com/jinzhu/gorm"
)

type Permission struct {
	gorm.Model
	RepoID  uint   `sql:"index;not null"` //fk
	UserID  uint   `sql:"index;not null"` //fk
	Type    string `sql:"not null"`       //commit / add branch / add tag / ...
	Branch  string `sql:"not null"`       //regex of branch
	Granted bool   `sql:"not null"`
}
