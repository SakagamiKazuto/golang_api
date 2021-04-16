package usecase

import "github.com/SakagamiKazuto/golang_api/domain"

type BosyuInteractor struct {
	Br BosyuRepository
}

func (bi BosyuInteractor) Create(bosyu *domain.Bosyu) (*domain.Bosyu, error) {
	return bi.Br.CreateBosyu(bosyu)
}

func (bi BosyuInteractor) GetByUid(uid uint) (*domain.Bosyus, error) {
	return bi.Br.FindBosyuByUid(uid)
}

func (bi BosyuInteractor) Update(bosyu *domain.Bosyu) (*domain.Bosyu, error) {
	return bi.Br.UpdateBosyu(bosyu)
}

func (bi BosyuInteractor) Delete(id uint) error {
	return bi.Br.DeleteBosyu(id)
}
