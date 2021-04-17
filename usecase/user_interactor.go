package usecase

import "github.com/SakagamiKazuto/golang_api/domain"

type UserInteractor struct {
	Ur UserRepository
}

func (ui UserInteractor) Create(user *domain.User) (*domain.User, error) {
	return ui.Ur.CreateUser(user)
}

func (ui UserInteractor) UserByMailPass (user *domain.User) (*domain.User,error) {
	return ui.Ur.FindUserByMailPass(user)
}

func (ui UserInteractor) UserByUid (user *domain.User) (*domain.User,error) {
	return ui.Ur.FindUserByUid(user)
}
