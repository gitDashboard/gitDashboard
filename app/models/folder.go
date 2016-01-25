package models

type Folder struct {
	ID           uint `gorm:"primary_key"`
	ParentID     uint
	Name         string `sql:"not null"`
	Path         string `sql:"not null"`
	Description  string
	Admins       []User `gorm:"many2many:admins_folders;"`
	Repositories []Repo
	Permissions  []Permission
}
