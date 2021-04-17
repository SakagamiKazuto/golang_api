package database

import "github.com/jinzhu/gorm"

type DBHandle interface {
	ConInf() *gorm.DB
}