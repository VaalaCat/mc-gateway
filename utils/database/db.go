package database

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load()
	initSqlite()
}

func GetDB() *gorm.DB {
	return GetSqlite()
}

// func CloseDB(db *gorm.DB) {
// 	tdb, err := db.DB()
// 	if err != nil {
// 		logrus.WithError(err).Errorf("Close DB error")
// 	}
// 	tdb.Close()
// }
