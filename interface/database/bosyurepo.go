package database

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/pkg/errors"
)

type BosyuRepository struct {
	DBHandle
}

func (br BosyuRepository) CreateBosyu(bosyu *domain.Bosyu) (*domain.Bosyu, error) {
	db := br.ConInf()
	if err := db.Create(&bosyu).Error; err != nil {
		return nil, CreateInDBError(err)
	}
	return bosyu, nil
}

func (br BosyuRepository) FindBosyuByUid(userID uint) (*domain.Bosyus, error) {
	db := br.ConInf()
	bosyus := &domain.Bosyus{}
	result := db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(bosyus)

	if result.Error != nil {
		return nil, CreateInDBError(result.Error)
	}
	// sliceでは.RecordNotFound()は使えない → https://qiita.com/hiromichi_n/items/a08a7e0f33641d71e6ef
	if result.RowsAffected == 0 {
		return nil, ExternalDBError{
			ErrorMessage:  fmt.Sprintf(`該当のユーザーID%dの募集は見つかりません`, userID),
			OriginalError: errors.New("record not found"),
			StatusCode:    ValueNotFound,
		}
	}
	return bosyus, nil
}

func (br BosyuRepository) UpdateBosyu(b *domain.Bosyu) (*domain.Bosyu, error) {
	db := br.ConInf()
	result := db.Model(b).Update(map[string]interface{}{
		"title":      b.Title,
		"about":      b.About,
		"prefecture": b.Prefecture,
		"city":       b.City,
		"level":      b.Level,
	})
	if result.Error != nil {
		return nil, CreateInDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ExternalDBError{
			ErrorMessage:  fmt.Sprintf(`該当の募集(ID=%d)は見つかりません`, b.ID),
			OriginalError: errors.New("record not found"),
			StatusCode:    ValueNotFound,
		}
	}

	return b, nil
}

func (br BosyuRepository) DeleteBosyu(bosyuID uint) error {
	db := br.ConInf()
	result := db.Where("id = ?", bosyuID).Delete(&domain.Bosyu{})
	if result.Error != nil {
		return CreateInDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ExternalDBError{
			ErrorMessage:  fmt.Sprintf(`該当の募集(ID=%d)は見つかりません`, bosyuID),
			OriginalError: errors.New("record not found"),
			StatusCode:    ValueNotFound,
		}
	}

	return nil
}
