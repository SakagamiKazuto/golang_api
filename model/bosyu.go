package model

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
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

func (b Bosyu) Validate() error {
	err := validation.ValidateStruct(&b,
		validation.Field(&b.Title, validation.Required),
		validation.Field(&b.About, validation.Required),
	)
	if err != nil {
		return &ExternalDBError{
			ErrorMessage:  fmt.Sprintf("タイトルまたは本文が空欄です"),
			OriginalError: err,
			StatusCode:    apperror.InvalidParameter,
		}
	}
	return nil
}

func CreateBosyu(bosyu *Bosyu, db *gorm.DB) (*Bosyu, error) {
	if err := db.Create(&bosyu).Error; err != nil {
		return nil, createInDBError(err)
	}
	return bosyu, nil
}

func FindBosyuByUid(userID uint, db *gorm.DB) ([]Bosyu, error) {
	var bosyus []Bosyu
	result := db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&bosyus)

	// sliceでは.RecordNotFound()は使えない → https://qiita.com/hiromichi_n/items/a08a7e0f33641d71e6ef
	if result.RowsAffected == 0 {
		return nil, ExternalDBError{
			ErrorMessage:  fmt.Sprintf(`該当のユーザーID%dの募集は見つかりません`, userID),
			OriginalError: errors.New("record not found"),
			StatusCode:    apperror.ValueNotFound,
		}
	}

	if result.Error != nil {
		return nil, createInDBError(result.Error)
	}
	return bosyus, nil
}

func UpdateBosyu(b *Bosyu, db *gorm.DB) (*Bosyu, error) {
	result := db.Model(b).Update(map[string]interface{}{
		"title":      b.Title,
		"about":      b.About,
		"prefecture": b.Prefecture,
		"city":       b.City,
		"level":      b.Level,
	})

	if result.RowsAffected == 0 {
		return nil, ExternalDBError{
			ErrorMessage:  fmt.Sprintf(`該当の募集(ID=%d)は見つかりません`, b.ID),
			OriginalError: errors.New("record not found"),
			StatusCode:    apperror.ValueNotFound,
		}
	}

	if result.Error != nil {
		return nil, createInDBError(result.Error)
	}

	return b, nil
}

func DeleteBosyu(bosyuID uint, db *gorm.DB) error {
	result := db.Where("id = ?", bosyuID).Delete(&Bosyu{})
	if result.RowsAffected == 0 {
		return ExternalDBError{
			ErrorMessage:  fmt.Sprintf(`該当の募集(ID=%d)は見つかりません`, bosyuID),
			OriginalError: errors.New("record not found"),
			StatusCode:    apperror.ValueNotFound,
		}
	}

	if result.Error != nil {
		return createInDBError(result.Error)
	}
	return nil
}
