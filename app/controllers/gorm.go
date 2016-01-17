package controllers

import (
	"database/sql"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
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
	event := &models.Event{}
	adminGrp := &models.Group{Name: "admin", Description: "Administration group"}
	if !Db.HasTable(group) {
		revel.INFO.Println("Creating groups table")
		Db.CreateTable(group)
		Db.Create(adminGrp)
	} else {
		Db.AutoMigrate(group)
	}
	if !Db.HasTable(user) {
		revel.INFO.Println("Creating users table")
		Db.CreateTable(user)
		//Db.Table("users_groups").AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT")
		//Db.Table("users_groups").AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")

		var adminPwdFld sql.NullString
		adminPwd, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		adminPwdFld.String = string(adminPwd)
		adminPwdFld.Valid = true
		adminUser := &models.User{Username: "admin", Type: "internal", Password: adminPwdFld, Name: "Admin user"}
		adminUser.Groups = []models.Group{*adminGrp}
		Db.Create(adminUser)
	} else {
		Db.AutoMigrate(user)
	}
	if !Db.HasTable(repo) {
		revel.INFO.Println("Creating repos table")
		Db.CreateTable(repo)
	} else {
		Db.AutoMigrate(repo)
	}
	if !Db.HasTable(perm) {
		revel.INFO.Println("Creating permissions table")
		Db.CreateTable(perm)
		Db.Model(perm).AddForeignKey("repo_id", "repos(id)", "CASCADE", "RESTRICT")
		Db.Model(perm).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
		Db.Model(perm).AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT")
	} else {
		Db.AutoMigrate(perm)
	}
	if !Db.HasTable(event) {
		revel.INFO.Println("Creating events table")
		Db.CreateTable(event)
	} else {
		Db.AutoMigrate(event)
	}
}

func (c *GormController) NewTransaction() *gorm.DB {
	txn := Db.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	revel.INFO.Println("txn init", txn)
	return txn
}

func (c *GormController) RollbackTransaction(tx *gorm.DB) {
	if tx == nil {
		return
	}
	tx.Rollback()
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	revel.INFO.Println("tx rollbacked", tx)
}
func (c *GormController) CommitTransaction(tx *gorm.DB) {
	if tx == nil {
		return
	}
	tx.Commit()
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	revel.INFO.Println("tx commited", tx)
}

func (c *GormController) Begin() revel.Result {
	c.Tx = c.NewTransaction()
	return nil
}
func (c *GormController) Commit() revel.Result {
	c.CommitTransaction(c.Tx)
	c.Tx = nil
	return nil
}

func (c *GormController) Rollback() revel.Result {
	c.RollbackTransaction(c.Tx)
	c.Tx = nil
	return nil
}
