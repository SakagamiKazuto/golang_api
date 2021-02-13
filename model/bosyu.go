package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Bosyu struct {
	gorm.Model
	Title      string    `json:"title"`
	About      string    `json:"about"`
	Prefecture string    `json:"pref"`
	City       string    `json:"city"`
	Level      string    `json:"level"`
	UserID     uint      `json:"user_id"`
	Message    []Message `gorm:"foreignkey:BosyuID"`
}

func CreateBosyu(b *Bosyu, db *gorm.DB) *Bosyu {
	db.Create(b)
	return b
}

func FindBosyu(user_id uint, db *gorm.DB) []Bosyu {
	var bosyus []Bosyu
	db.Where("user_id = ? AND deleted_at IS NULL", user_id).Find(&bosyus)
	return bosyus
}

func UpdateBosyu(b *Bosyu, db *gorm.DB) (*Bosyu, error) {
	rows := db.Model(b).Where("id = ?", b.ID).Update(map[string]interface{}{
		"title":      b.Title,
		"about":      b.About,
		"prefecture": b.Prefecture,
		"city":       b.City,
		"level":      b.Level,
	}).RowsAffected

	if rows == 0 {
		return nil, fmt.Errorf("Could not find Bosyu (ID = %v) in table", b.ID)
	}

	return b, nil
}

func DeleteBosyu(bosyu_id uint, db *gorm.DB) error {
	rows := db.Where("id = ?", bosyu_id).Delete(&Bosyu{}).RowsAffected
	if rows == 0 {
		return fmt.Errorf("Could not find Bosyu (ID = %v) in table", bosyu_id)
	}
	return nil
}


