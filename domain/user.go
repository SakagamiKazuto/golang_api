package domain

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Password string  `json:"password" gorm:"not null"`
	Name     string  `json:"name" gorm:"not null"`
	Address  string  `json:"address" `
	Tel      string  `json:"tel"`
	Mail     string  `json:"mail" gorm:"not null;unique"`
	Bosyu    []Bosyu `gorm:"foreignkey:UserID"`
}

type Users []User

