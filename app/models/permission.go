package models

import (
	"database/sql"
)

type Permission struct {
	ID       uint `gorm:"primary_key"`
	RepoID   uint `sql:"index;not null"` //fk
	UserID   sql.NullInt64
	GroupID  sql.NullInt64
	Type     string `sql:"not null"` //commit / add branch / add tag / ...
	Branch   string `sql:"not null"` //regex of branch
	Granted  bool   `sql:"not null"`
	Position uint   `sql:"not null"`
}
