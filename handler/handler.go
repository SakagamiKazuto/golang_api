package handler

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/apperror"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/SakagamiKazuto/golang_api/db"
	"github.com/SakagamiKazuto/golang_api/model"
)

// CreateBosyu is creating bosyu.
// @Summary create bosyu
// @Description create bosyu in a group
// @Accept  json
// @Produce  json
// @Success 201 {object} model.Bosyu
// @Failure 400 {object} echo.HTTPError
// @Router /api/bosyu/create [post]
func CreateBosyu(c echo.Context) error {
	if logined, err := isLogined(c); logined == false || err != nil {
		return createLoginFailureErr(c, err)
	}

	bosyu := new(model.Bosyu)
	if err := c.Bind(bosyu); err != nil {
		return apperror.ResponseError(c, ExternalHandleError{fmt.Sprintln(`リクエストに問題があります`), err, apperror.InvalidParameter})
	}


	if err := bosyu.Validate(); err != nil {
		return apperror.ResponseError(c, err)
	}

	bosyu, err := model.CreateBosyu(bosyu, db.DB)
	if err != nil {
		return apperror.ResponseError(c, err)
	}

	return c.JSON(http.StatusCreated, bosyu)
}

// GetBosyu is getting bosyu.
// @Summary get bosyu
// @Description get bosyu in a group
// @Accept  json
// @Produce  json
// @param user_id query string true "user_id which has bosyus"
// @Success 200 {object} []model.Bosyu
// @Failure 400,404 {object} echo.HTTPError
// @Router /api/bosyu/get [get]
func GetBosyu(c echo.Context) error {
	if logined, err := isLogined(c); logined == false || err != nil {
		return createLoginFailureErr(c, err)
	}

	bosyu := new(model.Bosyu)
	if err := c.Bind(bosyu); err != nil {
		return apperror.ResponseError(c, ExternalHandleError{err.Error(), err, apperror.InvalidParameter})
	}

	bosyus, err := model.FindBosyu(bosyu.UserID, db.DB)
	if err != nil {
		return apperror.ResponseError(c, err)
	}
	return c.JSON(http.StatusOK, bosyus)
}

// UpdateBosyu is updating bosyu.
// @Summary create bosyu
// @Description update bosyu in a group
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Bosyu
// @Failure 400,404 {object} echo.HTTPError
// @Router /api/bosyu/update [put]
func UpdateBosyu(c echo.Context) error {
	if logined, err := isLogined(c); logined == false || err != nil {
		return createLoginFailureErr(c, err)
	}

	bosyu := new(model.Bosyu)
	if err := c.Bind(bosyu); err != nil {
		return apperror.ResponseError(c, ExternalHandleError{err.Error(), err, apperror.InvalidParameter})
	}

	// validate必要

	b, err := model.UpdateBosyu(bosyu, db.DB)
	if err != nil {
		return apperror.ResponseError(c, err)
	}
	return c.JSON(http.StatusOK, b)
}

// DeleteBosyu is deleting bosyu.
// @Summary delete bosyu
// @Description delete bosyu in a group
// @Accept  json
// @Produce  json
// @param bosyu_id query string true "bosyu_id which bosyu has"
// @Success 200 {object} model.Bosyu
// @Failure 400,404 {object} echo.HTTPError
// @Router /api/bosyu/delete [delete]
func DeleteBosyu(c echo.Context) error {
	if logined, err := isLogined(c); logined == false || err != nil {
		return createLoginFailureErr(c, err)
	}

	bosyu := new(model.Bosyu)
	if err := c.Bind(bosyu); err != nil {
		return apperror.ResponseError(c, ExternalHandleError{err.Error(), err, apperror.InvalidParameter})
	}

	err := model.DeleteBosyu(bosyu.ID, db.DB)
	if err != nil {
		return apperror.ResponseError(c, ExternalHandleError{err.Error(), err, apperror.InvalidParameter})
	}
	return c.JSON(http.StatusOK, []string{fmt.Sprintf("募集(ID:%d)は削除されました", bosyu.ID)})
}
