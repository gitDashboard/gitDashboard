package controllers

import (
	"database/sql"
	"github.com/gitDashboard/gitDashboard/app/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"os"
)

var Db gorm.DB

type GormController struct {
	BasicController
	Tx *gorm.DB
}

func getDbReferences() (string, string) {
	var dbType, dbConnection string
	dbType, dbTypeFounded := revel.Config.String("db.type")
	dbConnection, dbConnectionFounded := revel.Config.String("db.connection")
	if !dbTypeFounded || !dbConnectionFounded {
		revel.WARN.Println("no settings found for database connection on revel configuration, searching on system environment (GITDASHBOARD_DBTYPE,GITDASHBOARD_DBCONNECTION)")
		//search on environment
		dbType = os.Getenv("GITDASHBOARD_DBTYPE")
		dbConnection = os.Getenv("GITDASHBOARD_DBCONNECTION")
	}
	if dbType == "" || dbConnection == "" {
		revel.ERROR.Println("No Database connection found")
		panic("No database connection found")
	}
	return dbType, dbConnection
}

func InitDB() {
	var err error
	dbtype, dbConnection := getDbReferences()
	Db, err = gorm.Open(dbtype, dbConnection)
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
	folder := &models.Folder{}
	event := &models.Event{}

	if !Db.HasTable(folder) {
		revel.INFO.Println("Creating folders table")
		Db.CreateTable(folder)
	} else {
		Db.AutoMigrate(folder)
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
		/*Db.Model(perm).AddForeignKey("repo_id", "repos(id)", "CASCADE", "RESTRICT")
		Db.Model(perm).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
		Db.Model(perm).AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT")*/
	} else {
		Db.AutoMigrate(perm)
	}
	if !Db.HasTable(user) {
		revel.INFO.Println("Creating users table")
		Db.CreateTable(user)
		Db.Table("users_permissions").AddForeignKey("permission_id", "permissions(id)", "CASCADE", "RESTRICT")
		Db.Table("users_permissions").AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")

		var adminPwdFld sql.NullString
		adminPwd, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		adminPwdFld.String = string(adminPwd)
		adminPwdFld.Valid = true
		adminUser := &models.User{Username: "admin", Type: "internal", Password: adminPwdFld, Name: "Admin user", Admin: true}
		Db.Create(adminUser)
	} else {
		Db.AutoMigrate(user)
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
