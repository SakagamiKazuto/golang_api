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

//// DeleteBosyu is deleting bosyu.
//// @Summary delete bosyu
//// @Description delete bosyu in a group
//// @Accept  json
//// @Produce  json
//// @param bosyu_id query string true "bosyu_id which bosyu has"
//// @Success 200 {object} model.Bosyu
//// @Failure 400,404 {object} echo.HTTPError
//// @Router /api/bosyu/delete [delete]
//func DeleteBosyu(c echo.Context) error {
//	if logined, err := isLogined(c); logined == false || err != nil {
//		return createLoginFailureErr(c, err)
//	}
//
//	bosyu := new(model.Bosyu)
//	if err := c.Bind(bosyu); err != nil {
//		return apperror.ResponseError(c, ExternalHandleError{err.Error(), err, apperror.InvalidParameter})
//	}
//
//	err := model.DeleteBosyu(bosyu.ID, db.DB)
//	if err != nil {
//		return apperror.ResponseError(c, ExternalHandleError{err.Error(), err, apperror.InvalidParameter})
//	}
//	return c.JSON(http.StatusOK, []string{fmt.Sprintf("募集(ID:%d)は削除されました", bosyu.ID)})
//}
