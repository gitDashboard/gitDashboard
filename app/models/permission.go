package models

import (
	"database/sql"
)

type Permission struct {
	ID       uint          `gorm:"primary_key"`
	RepoID   sql.NullInt64 `sql:"index"`    //fk -> repo
	FolderID sql.NullInt64 `sql:"index"`    //fk -> folder
	Type     string        `sql:"not null"` //commit / add branch / add tag / ...
	Branch   string        `sql:"not null"` //regex of branch
	Granted  bool          `sql:"not null"`
	Position uint          `sql:"not null"`
	Users    []User        `gorm:"many2many:users_permissions;"`
}
