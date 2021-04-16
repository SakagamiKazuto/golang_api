package usecase

import "github.com/SakagamiKazuto/golang_api/domain"

type BosyuRepository interface {
	CreateBosyu(*domain.Bosyu) (*domain.Bosyu, error)
	FindBosyuByUid(uint) (*domain.Bosyus, error)
	UpdateBosyu(*domain.Bosyu) (*domain.Bosyu, error)
	DeleteBosyu(uint) error
}
