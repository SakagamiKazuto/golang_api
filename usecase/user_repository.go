package usecase

import "github.com/SakagamiKazuto/golang_api/domain"

type UserRepository interface {
	CreateUser(*domain.User) (*domain.User, error)
	FindUserByMailPass(*domain.User) (*domain.User, error)
	FindUserByUid(*domain.User) (*domain.User, error)
}
