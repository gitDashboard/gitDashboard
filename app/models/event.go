package models

import (
	"time"
)

type Event struct {
	ID          uint `gorm:"primary_key"`
	RepoID      uint `sql:"index;not null"` //fk
	Reference   string
	Type        string `sql:"index;not null"`
	Level       string
	Started     time.Time  `sql:"DEFAULT:current_timestamp;not null"`
	Finished    *time.Time `sql:"index"`
	User        string
	Description string
}
