package database

import (
	"tg-mc/conf"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func initSqlite() {
	var err error
	godotenv.Load()

	dbPath := conf.GetBotSettings().DBPath
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logrus.Panic(err, "Initializing DB Error")
	}
	logrus.Info("Initialized DB at ", dbPath)
}

func GetSqlite() *gorm.DB {
	if db == nil {
		initSqlite()
	}
	return db
}
