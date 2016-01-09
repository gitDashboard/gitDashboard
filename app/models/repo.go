package models

type Repo struct {
	ID          uint   `gorm:"primary_key"`
	Path        string `sql:"not null"`
	Permissions []Permission
}
