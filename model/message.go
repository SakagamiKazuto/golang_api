package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Message struct {
	gorm.Model
	Message     string `json:"message"`
	BosyuID uint `json:"bosyu_id"`
	User User
}

//func CreateUser(user *User) {
//	db.Create(user)
//}
//
//func FindUser(u *User) User {
//	var user User
//	db.Where(u).First(&user)
//	return user
//}