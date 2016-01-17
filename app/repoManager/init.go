package repoManager

import (
	"errors"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/jinzhu/gorm"

	"time"
)

func GetRepo(db *gorm.DB, fullPath string) (models.Repo, error) {
	var repo models.Repo
	dbExec := db.Where("path = ? ", fullPath).First(&repo)
	if dbExec.Error != nil {
		return repo, dbExec.Error
	}
	return repo, nil
}

func AddRepoEvent(db *gorm.DB, repoId uint, eventType, user, description string) (*models.Event, error) {
	now := time.Now()
	dbEvent := models.Event{RepoID: repoId, Type: eventType, User: user, Description: description, Started: now, Finished: &now}
	dbEx := db.Create(&dbEvent)
	if dbEx.Error != nil {
		return nil, dbEx.Error
	}
	return &dbEvent, nil
}

func StartRepoEvent(db *gorm.DB, repoId uint, eventType, user, description string) (*models.Event, error) {
	now := time.Now()
	dbEvent := models.Event{RepoID: repoId, Type: eventType, User: user, Description: description, Started: now}
	dbEx := db.Create(&dbEvent)
	if dbEx.Error != nil {
		return nil, dbEx.Error
	}
	return &dbEvent, nil
}

func FinishRepoEvent(db *gorm.DB, eventId uint) error {
	now := time.Now()
	var dbEvent models.Event
	dbFind := db.First(&dbEvent, eventId)
	if dbFind.Error != nil {
		return dbFind.Error
	}
	if dbEvent.ID != eventId {
		return errors.New("Event not found")
	}

	dbEvent.Finished = &now
	dbEx := db.Save(&dbEvent)
	if dbEx.Error != nil {
		return dbEx.Error
	}
	return nil
}

func HasOperationInProgress(db *gorm.DB, repoId uint) (bool, error) {
	var events []models.Event
	dbFind := db.Where("repo_id = ? AND finished IS NULL", repoId).Find(&events)
	if dbFind.Error != nil {
		return false, dbFind.Error
	} else {
		return len(events) > 0, nil
	}
}
