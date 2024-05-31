package models

import (
	"tg-mc/utils/database"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TGID   int64  `gorm:"column:tgid;NOT NULL;index"`
	MCName string `gorm:"column:mc_name;NOT NULL;unique"`
	Status int    `gorm:"column:status;NOT NULL"` // 0: pending, 1: normal, 2: banned
}

func init() {
	if err := database.GetDB().AutoMigrate(&User{}); err != nil {
		logrus.Panic(err)
	}
}

func (u *User) TableName() string {
	return "users"
}

func GetUsersByTGID(tgID int64) (users []User, err error) {
	err = database.GetDB().Where(
		&User{TGID: tgID},
	).Find(&users).Error
	return
}

func GetUserByTGID(tgID int64) (user User, err error) {
	err = database.GetDB().Where(
		&User{TGID: tgID},
	).First(&user).Error
	return
}

func GetUserByMCName(mcName string) (user User, err error) {
	err = database.GetDB().Where(
		&User{MCName: mcName},
	).First(&user).Error
	return
}

func CreateUser(u *User) (err error) {
	err = database.GetDB().Create(&u).Error
	return
}

func (u *User) Delete(tgID int64) error {
	return database.GetDB().Where(
		&User{TGID: tgID},
	).Unscoped().Delete(&u).Error
}
