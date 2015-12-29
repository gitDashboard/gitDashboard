package controllers

import (
	"crypto/md5"
	"database/sql"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/revel/revel"
)

var Db gorm.DB

type GormController struct {
	BasicController
	Tx *gorm.DB
}

func InitDB() {
	var err error
	Db, err = gorm.Open("postgres", "user=igor password=infocam dbname=gitdashboard sslmode=disable")
	if err != nil {
		revel.ERROR.Println("FATAL", err)
		panic(err)
	}
	Db.LogMode(true)
	Db.DB().Ping()
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)

	perm := &models.Permission{}
	user := &models.User{}
	repo := &models.Repo{}
	group := &models.Group{}
	adminGrp := &models.Group{Name: "admin", Description: "Administration group"}
	if !Db.HasTable(group) {
		revel.INFO.Println("Creating groups table")
		Db.CreateTable(group)
		Db.Create(adminGrp)
	}
	if !Db.HasTable(user) {
		revel.INFO.Println("Creating users table")
		Db.CreateTable(user)
		Db.Table("users_groups").AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT")
		Db.Table("users_groups").AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")

		var adminPwdFld sql.NullString
		adminPwd := md5.Sum([]byte("admin"))
		adminPwdFld.String = string(adminPwd[:16])
		adminUser := &models.User{Username: "admin", Type: "internal", Password: adminPwdFld}
		adminUser.Groups = []models.Group{*adminGrp}
		Db.Create(adminUser)
	}
	if !Db.HasTable(repo) {
		revel.INFO.Println("Creating repos table")
		Db.CreateTable(repo)
	}
	if !Db.HasTable(perm) {
		revel.INFO.Println("Creating permissions table")
		Db.CreateTable(perm)
		Db.Model(perm).AddForeignKey("repo_id", "repos(id)", "CASCADE", "RESTRICT")
		Db.Model(perm).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
		Db.Model(perm).AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT")
	}
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
