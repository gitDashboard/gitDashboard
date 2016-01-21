package models

type Folder struct {
	ID           uint `gorm:"primary_key"`
	ParentID     uint
	Name         string `sql:"not null"`
	Path         string `sql:"not null"`
	Description  string
	Repositories []Repo
	Permissions  []Permission
}
