package apperror

import (
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ErrorResponse struct {
	Code   database.ErrorCode `json:"code"`
	Errors []string           `json:"errors"`
}

func ResponseError(c echo.Context, err error) error {
	exe, ok := errors.Cause(err).(database.ExternalError)
	if ok {
		log.WithFields(log.Fields{
			"ErrCode":    exe.Code(),
			"ErrMessage": err.Error(),
		}).Info("External Error occurred.")

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
		log.WithFields(log.Fields{
			"ErrCode":    database.InternalServerError,
			"ErrMessage": err.Error(),
		}).Warn("Internal Error occurred.")
		return c.JSON(http.StatusInternalServerError, ieres)
	}

	// Unhandledなエラーの処理
	ieres.Code = database.UnHandledError
	log.WithFields(log.Fields{
		"ErrCode":    database.UnHandledError,
		"ErrMessage": err.Error(),
	}).Warn("Unhandled Error occurred.")
	return c.JSON(http.StatusInternalServerError, ieres)
}

