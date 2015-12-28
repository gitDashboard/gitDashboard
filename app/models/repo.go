package models

import (
	"github.com/jinzhu/gorm"
)

type Repo struct {
	gorm.Model
	Path        string `sql:"not null"`
	Permissions []Permission
}
