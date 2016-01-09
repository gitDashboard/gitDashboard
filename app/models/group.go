package models

type Group struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `sql:"not null;unique_index"`
	Description string
	Users       []User `gorm:"many2many:users_groups;"`
}
