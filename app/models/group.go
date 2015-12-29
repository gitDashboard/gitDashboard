package models

import (
	"github.com/jinzhu/gorm"
)

type Group struct {
	gorm.Model
	Name        string `sql:"not null;unique_index"`
	Description string
}
