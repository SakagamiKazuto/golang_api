package apperror

import (
	l "github.com/SakagamiKazuto/golang_api/infra/waf/logger"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type ErrorResponse struct {
	Code   database.ErrorCode `json:"code"`
	Errors []string           `json:"errors"`
}

func ResponseError(c echo.Context, err error) error {
	exe, ok := errors.Cause(err).(database.ExternalError)
	if ok {
		l.Log.InfoWithFields("External Error occurred.", database.Fields{
			"ErrCode":    exe.Code(),
			"ErrMessage": err.Error(),
		})

		httpStatus := database.GetHttpStatus(exe.Code())
		return c.JSON(httpStatus, &ErrorResponse{
			Code:   exe.Code(),
			Errors: exe.Messages(),
		})
	}

	ieres := &ErrorResponse{Errors: []string{"サーバー内部でエラーが発生しました"}}

	ie, ok := errors.Cause(err).(database.InternalError)
	// handledなエラーの処理
	if ok && ie.Internal() {
		ieres.Code = database.InternalServerError
		l.Log.WarnWithFields("Internal Error occurred.", database.Fields{
			"ErrCode":    database.InternalServerError,
			"ErrMessage": err.Error(),
		})
		return c.JSON(http.StatusInternalServerError, ieres)
	}

	// Unhandledなエラーの処理
	ieres.Code = database.UnHandledError
	l.Log.WarnWithFields("Unhandled Error occurred.", database.Fields{
		"ErrCode":    database.UnHandledError,
		"ErrMessage": err.Error(),
	})
	return c.JSON(http.StatusInternalServerError, ieres)
}
