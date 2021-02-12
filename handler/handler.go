package handler

import (
	//"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"work/db"
	"work/model"
)

/* Tasks
!!JWT入れた際の認証のロジック追加
!!BosyuGetに関してparamがblankだった場合、勝手に0が入る件解決
*/

// CreateBosyu is creating bosyu.
// @Summary create bosyu
// @Description create bosyu in a group
// @Accept  json
// @Produce  json
// @Success 201 {object} model.Bosyu
// @Failure 400 {object} echo.HTTPError
// @Router /api/bosyu/create [post]
func CreateBosyu(c echo.Context) error {
	if CheckHasLogined(c) == false {
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

func GetBosyu(c echo.Context) error {
	//JWTを入れた際に必要になるロジック
	if CheckHasLogined(c) == false {
		return &echo.HTTPError{
			Code: http.StatusNotFound,
			Message: "can't find login user.",
		}
	}

	if c.QueryParam("user_id") == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Blank user_id is invalid parameter",
		}
	}

	user_id, _ := strconv.ParseUint(c.QueryParam("user_id"), 10, 32)

	bosyus := model.FindBosyu(uint(user_id), db.DB)
	if len(bosyus) == 0 {
		return &echo.HTTPError{
			Code: http.StatusNotFound,
			Message: "can't find bosyus.",
		}
	}
	return c.JSON(http.StatusOK, bosyus)
}

//func UpdateBosyu(c echo.Context) error {
//	// JWTを入れた際に必要になるロジック追加
//	bosyu := new(model.Bosyu)
//	if err := c.Bind(bosyu); err != nil {
//		return err
//	}
//
//	if bosyu.Title == "" || bosyu.About == "" {
//		return &echo.HTTPError{
//			Code:    http.StatusBadRequest,
//			Message: "invalid to Title or About fields",
//		}
//	}
//
//	model.UpdateBosyu(bosyu)
//
//	return c.JSON(http.StatusCreated, bosyu)
//}
//
//func DeleteBosyu(c echo.Context) error {
//	// JWTを入れた際にはロジック追加
//	bosyu_id, err := strconv.ParseUint(c.QueryParam("id"), 10, 32)
//
//	if err != nil {
//		return echo.ErrNotFound
//	}
//
//	model.DeleteBosyu(uint(bosyu_id))
//
//	return c.NoContent(http.StatusNoContent)
//}

func CheckHasLogined(c echo.Context) bool {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{Model: gorm.Model{ID: uid}}, db.DB); user.ID == 0 {
		return false
	}
	return true
}