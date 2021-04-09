package apperror

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ErrorResponse struct {
	Code   ErrorCode `json:"code"`
	Errors []string  `json:"errors"`
}

type externalError interface {
	Code() ErrorCode
	Messages() []string
}

type internalError interface {
	Internal() bool
}

type ErrorCode int

const (
	AuthenticationParamMissing ErrorCode = iota
	AuthenticationFailure
	InvalidParameter
	UniqueValueDuplication
	ValueNotFound
	InternalError

	// Error codes for unhandled error
	UnHandledError ErrorCode = 999
)

var codeStatusMap = map[ErrorCode]int{
	AuthenticationFailure:      http.StatusForbidden,
	AuthenticationParamMissing: http.StatusBadRequest,
	InvalidParameter:           http.StatusBadRequest,
	UniqueValueDuplication:     http.StatusBadRequest,
	ValueNotFound:              http.StatusNotFound,
	InternalError:              http.StatusInternalServerError,
}

func ResponseError(ec echo.Context, err error) error {
	exe, ok := errors.Cause(err).(externalError)
	if ok {
		log.WithFields(log.Fields{
			"ErrCode":    exe.Code(),
			"ErrMessage": err.Error(),
		}).Info("External Error occurred.")

		httpStatus := GetHttpStatus(exe.Code())
		return ec.JSON(httpStatus, &ErrorResponse{
			Code:   exe.Code(),
			Errors: exe.Messages(),
		})
	}

	ieres := &ErrorResponse{Errors: []string{"サーバー内部でエラーが発生しました"}}

	ie, ok := errors.Cause(err).(internalError)
	// handledなエラーの処理
	if ok && ie.Internal() {
		ieres.Code = InternalError
		log.WithFields(log.Fields{
			"ErrCode":    InternalError,
			"ErrMessage": err.Error(),
		}).Warn("Internal Error occurred.")
		return ec.JSON(http.StatusInternalServerError, ieres)
	}

	// Unhandledなエラーの処理
	ieres.Code = UnHandledError
	log.WithFields(log.Fields{
		"ErrCode":    UnHandledError,
		"ErrMessage": err.Error(),
	}).Warn("Unhandled Error occurred.")
	return ec.JSON(http.StatusInternalServerError, ieres)
}

func GetHttpStatus(code ErrorCode) int {
	return codeStatusMap[code]
}
