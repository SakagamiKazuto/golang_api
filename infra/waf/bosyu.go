package waf

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/SakagamiKazuto/golang_api/infra/waf/apperror"
	"github.com/SakagamiKazuto/golang_api/interface/controller"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

//// CreateBosyu is creating bosyu.
//// @Summary create bosyu
//// @Description create bosyu in a group
//// @Accept  json
//// @Produce  json
//// @Success 201 {object} model.Bosyu
//// @Failure 400 {object} echo.HTTPError
//// @Router /api/bosyu/create [post]
func (s Server) CreateBosyu(bc *controller.BosyuController) echo.HandlerFunc {
	return C(func(c *Context) error {
		bosyu := new(domain.Bosyu)
		if err := c.Bind(bosyu); err != nil {
			return apperror.ResponseError(c, database.ExternalDBError{fmt.Sprintln(`リクエストに問題があります`), err, database.InvalidParameter})
		}
		if err := ValidateBosyu(bosyu); err != nil {
			return apperror.ResponseError(c, err)
		}
		b, err := bc.CreateBosyu(bosyu)
		if err != nil {
			return apperror.ResponseError(c, err)
		}
		return c.JSON(http.StatusCreated, b)
	})
}

//// GetBosyu is getting bosyu.
//// @Summary get bosyu
//// @Description get bosyu in a group
//// @Accept  json
//// @Produce  json
//// @param user_id query string true "user_id which has bosyus"
//// @Success 200 {object} []model.Bosyu
//// @Failure 400,404 {object} echo.HTTPError
//// @Router /api/bosyu/get [get]
func (s Server) GetBosyu(bc *controller.BosyuController) echo.HandlerFunc {
	return C(func(c *Context) error {
		bosyu := new(domain.Bosyu)
		if err := c.Bind(bosyu); err != nil {
			return apperror.ResponseError(c, database.ExternalDBError{err.Error(), err, database.InvalidParameter})
		}
		bosyus, err := bc.FindBosyuByUid(bosyu.UserID)
		if err != nil {
			return apperror.ResponseError(c, err)
		}
		return c.JSON(http.StatusOK, bosyus)
	})
}

//// UpdateBosyu is updating bosyu.
//// @Summary create bosyu
//// @Description update bosyu in a group
//// @Accept  json
//// @Produce  json
//// @Success 200 {object} model.Bosyu
//// @Failure 400,404 {object} echo.HTTPError
//// @Router /api/bosyu/update [put]
func (s Server) UpdateBosyu(bc *controller.BosyuController) echo.HandlerFunc {
	return C(func(c *Context) error {
		bosyu := new(domain.Bosyu)
		if err := c.Bind(bosyu); err != nil {
			return apperror.ResponseError(c, database.ExternalDBError{err.Error(), err, database.InvalidParameter})
		}
		if err := ValidateBosyu(bosyu); err != nil {
			return apperror.ResponseError(c, err)
		}

		b, err := bc.UpdateBosyu(bosyu)
		if err != nil {
			return apperror.ResponseError(c, err)
		}
		return c.JSON(http.StatusOK, b)
	})
}

func (s Server) DeleteBosyu(bc *controller.BosyuController) echo.HandlerFunc {
	return C(func(c *Context) error {
		bosyu := new(domain.Bosyu)
		if err := c.Bind(bosyu); err != nil {
			return apperror.ResponseError(c, database.ExternalDBError{err.Error(), err, database.InvalidParameter})
		}
		err :=bc.DeleteBosyu(bosyu.ID)
		if err != nil {
			return apperror.ResponseError(c, err)
		}
		return c.JSON(http.StatusOK, []string{fmt.Sprintf("募集(ID:%d)は削除されました", bosyu.ID)})
	})
}
