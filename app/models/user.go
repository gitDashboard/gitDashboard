package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username    string `sql:"not null;index"`
	Type        string //internal/LDAP/...
	Password    sql.NullString
	Permissions []Permission
}
