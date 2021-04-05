package model

import (
	"github.com/SakagamiKazuto/golang_api/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (u User) Validate() error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Mail, validation.Required),
	)
	if err != nil {
		return &apperror.ExternalError{
			ErrorMessage:  err.Error(),
			OriginalError: err,
			StatusCode:    apperror.InvalidParameter,
		}
	}
	return nil
}

// !! external2つと、internal1つDBでエラーハンドルする
func CreateUser(user *User, db *gorm.DB) (*User, error) {
	err := db.Create(user).Error
	if err != nil {
		return nil, &apperror.ExternalError{
			ErrorMessage:  err.Error(),
			OriginalError: err,
			StatusCode:    apperror.UniqueValueDuplication,
		}
	}

	return user, nil
}

func FindUser(u *User, db *gorm.DB) User {
	var user User
	db.Where(u).First(&user)
	return user
}

// external, internalそれぞれの構造体を返すように
//type ExternalError struct {
//	errorMessage string
//	originalError error
//	statusCode int
//}
//
//func (de ExternalError) Messages() []string {
//	return []string{de.errorMessage}
//}
//
//func (de ExternalError) Code() apperror.ErrorCode {
//	return apperror.ErrorCode(de.statusCode)
//}
//
//func (de ExternalError) Error() string {
//	return de.originalError.Error()
//}
