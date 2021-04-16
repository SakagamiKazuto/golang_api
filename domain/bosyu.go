package domain

import (
	"github.com/jinzhu/gorm"
)

type Bosyu struct {
	gorm.Model
	Title      string    `json:"title" gorm:"not null"`
	About      string    `json:"about" gorm:"not null"`
	Prefecture string    `json:"pref"`
	City       string    `json:"city"`
	Level      string    `json:"level"`
	UserID     uint      `json:"user_id"`
	Message    []Message `gorm:"foreignkey:BosyuID"`
}

type Bosyus []Bosyu

