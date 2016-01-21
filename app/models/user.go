package models

import (
	"database/sql"
)

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `sql:"not null;index"`
	Type     string //internal/LDAP/...
	Password sql.NullString
	Name     string `sql:"not null"`
	Email    sql.NullString
	Admin    bool
}
