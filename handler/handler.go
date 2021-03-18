package handler

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

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
	if checkHasLogined(c) == false {
		return &echo.HTTPError{
			Code: http.StatusNotFound,
			Message: "can't find login user.",
		}
	}

	bosyu := new(model.Bosyu)
	if err := c.Bind(bosyu); err != nil {
		return err
	}

	if bosyu.Title == "" || bosyu.About == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid to Title or About fields",
		}
	}

	bosyu = model.CreateBosyu(bosyu, db.DB)

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
	if checkHasLogined(c) == false {
		return &echo.HTTPError{
			Code: http.StatusNotFound,
			Message: "can't find login user.",
		}
	}


	user_id, err := strconv.ParseUint(c.QueryParam("user_id"), 10, 32)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid user_id parameter",
		}
	}

	bosyus := model.FindBosyu(uint(user_id), db.DB)
	if len(bosyus) == 0 {
		return &echo.HTTPError{
			Code: http.StatusNotFound,
			Message: "can't find bosyus.",
		}
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
	if checkHasLogined(c) == false {
		return &echo.HTTPError{
			Code: http.StatusNotFound,
			Message: "can't find login user.",
		}
	}

	bosyu := new(model.Bosyu)
	var err error
	if err = c.Bind(bosyu); err != nil {
		return err
	}

	if bosyu.Title == "" || bosyu.About == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid to Title or About fields",
		}
	}

	bosyu, err = model.UpdateBosyu(bosyu, db.DB)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}
	}
	return c.JSON(http.StatusOK, bosyu)
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
	if checkHasLogined(c) == false {
		return &echo.HTTPError{
			Code: http.StatusNotFound,
			Message: "can't find login user.",
		}
	}

	bosyu_id, err := strconv.ParseUint(c.QueryParam("bosyu_id"), 10, 32)

	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid bosyu id parameter",
		}
	}

	err = model.DeleteBosyu(uint(bosyu_id), db.DB)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func checkHasLogined(c echo.Context) bool {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{Model: gorm.Model{ID: uid}}, db.DB); user.ID == 0 {
		return false
	}
	return true
}