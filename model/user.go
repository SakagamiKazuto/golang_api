package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"work/db"
)

type User struct {
	gorm.Model
	Password string  `json:"password"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Tel      string  `json:"tel"`
	Mail     string  `json:"mail"`
	Bosyu    []Bosyu `gorm:"foreignkey:UserID"`
}

func CreateUser(user *User, db *gorm.DB) {
	db.Create(user)
}

// !!パスワードも使ってログインするように実装を変更したい
func FindUser(u *User, db *gorm.DB) User {
	var user User
	// SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1
	db.Where(u).First(&user)
	return user
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
