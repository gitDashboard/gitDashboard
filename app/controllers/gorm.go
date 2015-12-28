package controllers

import (
	"database/sql"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/revel/revel"
)

var Db gorm.DB

type GormController struct {
	*revel.Controller
	Tx *gorm.DB
}

func InitDB() {
	var err error
	Db, err = gorm.Open("postgres", "user=igor password=infocam dbname=gitdashboard sslmode=disable")
	if err != nil {
		revel.ERROR.Println("FATAL", err)
		panic(err)
	}
	Db.DB().Ping()
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)

	perm := &models.Permission{}
	user := &models.User{}
	repo := &models.Repo{}
	if !Db.HasTable(user) {
		Db.CreateTable(user)
	}
	if !Db.HasTable(repo) {
		Db.CreateTable(user, repo)
	}
	if !Db.HasTable(perm) {
		Db.CreateTable(perm)
		Db.Model(perm).AddForeignKey("repo_id", "repos(id)", "CASCADE", "RESTRICT")
		Db.Model(perm).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	}
	Db.LogMode(true)

}

func (c *GormController) Begin() revel.Result {
	txn := Db.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	c.Tx = txn
	revel.INFO.Println("c.Tx init", c.Tx)
	return nil
}
func (c *GormController) Commit() revel.Result {
	if c.Tx == nil {
		return nil
	}
	c.Tx.Commit()
	if err := c.Tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Tx = nil
	revel.INFO.Println("c.Tx commited (nil)")
	return nil
}

func (c *GormController) Rollback() revel.Result {
	if c.Tx == nil {
		return nil
	}
	c.Tx.Rollback()
	if err := c.Tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Tx = nil
	return nil
}
