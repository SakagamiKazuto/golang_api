package controller

import (
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"github.com/SakagamiKazuto/golang_api/usecase"
)

type UserController struct {
	Itr usecase.UserInteractor
}

func NewUserController(dbHandle database.DBHandle) *UserController {
	return &UserController{
		Itr: usecase.UserInteractor{
			Ur: database.UserRepository{
				dbHandle,
			},
		},
	}
}


func (uc UserController) CreateUser(user *domain.User) (*domain.User, error) {
	return uc.Itr.Create(user)
}

func (uc UserController) GetUserByMailPass(user *domain.User) (*domain.User, error) {
	return uc.Itr.UserByMailPass(user)
}
