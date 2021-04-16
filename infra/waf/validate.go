package waf

import (
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateUser(u *domain.User) error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Mail, validation.Required),
	)
	if err != nil {
		return &database.ExternalDBError{
			ErrorMessage:  "バリデーションに失敗しました",
			OriginalError: err,
			StatusCode:    database.InvalidParameter,
		}
	}
	return nil
}
func ValidateBosyu(b *domain.Bosyu) error {
	err := validation.ValidateStruct(b,
		validation.Field(&b.Title, validation.Required),
		validation.Field(&b.About, validation.Required),
	)
	if err != nil {
		return &database.ExternalDBError{
			ErrorMessage:  "バリデーションに失敗しました",
			OriginalError: err,
			StatusCode:    database.InvalidParameter,
		}
	}
	return nil
}
