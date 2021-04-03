package model

import (
	"github.com/SakagamiKazuto/golang_api/apperror"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
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

func CreateUser(user *User, db *gorm.DB) (*User, error) {
	if user.Name == "" || user.Mail == "" || user.Password == "" {
		err := errors.New("Name, Mail, Password must not be zero-value")
		return nil, &apperror.ExternalError{
			ErrorMessage: err.Error(),
			OriginalError: err,
			StatusCode: apperror.AuthenticationParamMissing,
		}
	}

	err := db.Create(user).Error
	if err != nil {
		return nil, &apperror.ExternalError{
			ErrorMessage: err.Error(),
			OriginalError: err,
			StatusCode: apperror.UniqueValueDuplication,
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
