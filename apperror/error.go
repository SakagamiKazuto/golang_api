package apperror

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ErrorContext struct {
	echo.Context
}

func (ec *ErrorContext) ResponseError(err error) error {
	exe, ok := errors.Cause(err).(externalError)
	if ok {
		log.WithFields(log.Fields{
			"ErrCode": exe.Code(),
			"ErrMessage": err.Error(),
		}).Info("External Error occurred.")

		httpStatus := GetHttpStatus(exe.Code())
		return ec.JSON(httpStatus, &ErrorResponse{
			Code:   exe.Code(),
			Errors: exe.Messages(),
		})
	}

	ieres := &ErrorResponse{Errors: []string{err.Error()}}

	ie, ok := errors.Cause(err).(internalError)
	// Unhandledなエラーの処理
	if !ok {
		ieres.Code = UnHandledError

		log.WithFields(log.Fields{
			"ErrCode": UnHandledError,
			"ErrMessage": err.Error(),
		}).Warn("Unhandled Error occurred.")
		return ec.JSON(http.StatusInternalServerError, ieres)
	}

	// handledなエラーの処理
	if ie.Internal() {
		ieres.Code = InternalError
		//!! message, 発生場所を出力。 messageはinternalErrorの実装に依存
		log.WithFields(log.Fields{
			"ErrCode": InternalError,
			"ErrMessage": err.Error(),
		}).Warn("Internal Error occurred.")
		return ec.JSON(http.StatusInternalServerError, ieres)
	}

	return ec.JSON(http.StatusInternalServerError, ieres)
}

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
	InternalError

	// Error codes for internal error
	UnHandledError ErrorCode = 999
)

var codeStatusMap = map[ErrorCode]int{
	AuthenticationFailure:      http.StatusForbidden,
	AuthenticationParamMissing: http.StatusBadRequest,
	InvalidParameter:           http.StatusBadRequest,
	UniqueValueDuplication:     http.StatusBadRequest,
	InternalError:              http.StatusInternalServerError,
}

func GetHttpStatus(code ErrorCode) int {
	return codeStatusMap[code]
}

