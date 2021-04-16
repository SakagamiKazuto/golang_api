package database

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/lib/pq"
)

type UserRepository struct {
	DBHandle
}

func (ur UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	db := ur.ConInf()
	if err := db.Create(&user).Error; err != nil {
		pqe, ok := err.(*pq.Error)
		if !ok {
			return nil, err
		}

		switch pqe.Code {
		//DB内でUniqueKey制約に引っかかるエラーの場合にはexternalエラーを返す
		case "23505":
			return nil, &ExternalDBError{
				ErrorMessage:  fmt.Sprintf("メールアドレス%sのデータ挿入に失敗しました", user.Mail),
				OriginalError: pqe,
				StatusCode:    UniqueValueDuplication,
			}
		default:
			return nil, &InternalDBError{
				Message:       pqe.Message,
				Detail:        pqe.Detail,
				File:          pqe.File,
				Line:          pqe.Line,
				OriginalError: pqe,
			}
		}
	}

	return user, nil
}

func (ur UserRepository) FindUserByMailPass(u *domain.User) (*domain.User, error) {
	db := ur.ConInf()
	user := new(domain.User)
	result := db.Where("mail = ? AND password = ?", u.Mail, u.Password).First(&user)

	if result.RecordNotFound() {
		return nil, ExternalDBError{
			ErrorMessage:  fmt.Sprintln("該当のユーザーが見つかりません"),
			OriginalError: result.Error,
			StatusCode:    ValueNotFound,
		}
	}

	if result.Error != nil {
		return nil, CreateInDBError(result.Error)
	}
	return user, nil
}

func (ur UserRepository) FindUserByUid(u *domain.User) (*domain.User, error) {
	db := ur.ConInf()
	user := new(domain.User)
	result := db.Where("id = ?", u.ID).First(&user)

	if result.RecordNotFound() {
		return nil, ExternalDBError{
			ErrorMessage:  fmt.Sprintln("該当のユーザーが見つかりません"),
			OriginalError: result.Error,
			StatusCode:    ValueNotFound,
		}
	}

	if result.Error != nil {
		return nil, CreateInDBError(result.Error)
	}
	return user, nil
}
