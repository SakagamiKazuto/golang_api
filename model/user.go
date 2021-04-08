package model

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lib/pq"
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
		return &ExternalDBError{
			ErrorMessage:  err.Error(),
			OriginalError: err,
			StatusCode:    apperror.InvalidParameter,
		}
	}
	return nil
}

func CreateUser(user *User, db *gorm.DB) (*User, error) {
	if err := db.Create(&user).Error; err != nil {
		pqe, ok := err.(*pq.Error)
		if !ok {
			return nil, err
		}


		switch pqe.Code {
		// DB内でUniqueKey制約に引っかかるエラーの場合にはexternalエラーを返す
		case "23505":
			return nil, &ExternalDBError{
				ErrorMessage:  fmt.Sprintf("メールアドレス%sのデータ挿入に失敗しました", user.Mail),
				OriginalError: pqe,
				StatusCode:    apperror.UniqueValueDuplication,
			}
		default:
			return nil, &InternalDBError{
				ErrorMessage:  createInErrMsg(1),
				OriginalError: pqe,
			}
		}
	}

	return user, nil
}

func FindUser(u *User, db *gorm.DB) User {
	var user User
	db.Where(u).First(&user)
	return user
}

