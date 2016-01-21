package models

type Repo struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `sql:"not null"`
	Path        string `sql:"not null;unique"`
	Locked      bool
	FolderID    uint //fk -> Folder
	Permissions []Permission
	Events      []Event
}
