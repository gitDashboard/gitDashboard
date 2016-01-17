package models

import (
	"time"
)

type Event struct {
	ID          uint      `gorm:"primary_key"`
	RepoID      uint      `sql:"index;not null"` //fk
	Type        string    `sql:"index;not null"`
	Started     time.Time `sql:"DEFAULT:current_timestamp;not null"`
	Finished    *time.Time
	User        string
	Description string
}
