package controller

import (
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"github.com/SakagamiKazuto/golang_api/usecase"
)

type BosyuController struct {
	Itr usecase.BosyuInteractor
}

func NewBosyuController(dbHandle database.DBHandle) *BosyuController {
	return &BosyuController{
		Itr: usecase.BosyuInteractor{
			Br: database.BosyuRepository{
				dbHandle,
			},
		},
	}
}

func (bc BosyuController) CreateBosyu(bosyu *domain.Bosyu) (*domain.Bosyu, error) {
	return bc.Itr.Create(bosyu)
}

func (bc BosyuController) FindBosyuByUid(uid uint) (*domain.Bosyus, error) {
	return bc.Itr.GetByUid(uid)
}

func (bc BosyuController) UpdateBosyu(bosyu *domain.Bosyu) (*domain.Bosyu, error) {
	return bc.Itr.Update(bosyu)
}

func (bc BosyuController) DeleteBosyu(id uint) error {
	return bc.Itr.Delete(id)
}

