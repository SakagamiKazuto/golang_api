package apperror

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type ErrorContext struct {
	echo.Context
}

func (ec *ErrorContext) ResponseError(err error) error {
	exe, ok := errors.Cause(err).(externalError)
	if ok {
		httpStatus := GetHttpStatus(exe.Code())
		eres := &ErrorResponse{
			Code:   exe.Code(),
			Errors: exe.Messages(),
		}
		return ec.JSON(httpStatus, eres)
	}

	ie, ok := errors.Cause(err).(internalError)
	if ok && ie.Internal() {
		eres := &ErrorResponse{
			Code:   InternalError,
			Errors: []string{"We are very sorry, internal error occurred. We will start investigation immediately."},
		}
		// ここにinternal error用のloggerを設定する
		return ec.JSON(http.StatusInternalServerError, eres)
	}

	// ここにunhandled error用のloggerを設定する
	return ec.JSON(http.StatusInternalServerError, &ErrorResponse{
		Code:   InternalError,
		Errors: []string{"Unexpected Error"},
	})
}

type ErrorResponse struct {
	//Code             internal.ErrorCode `json:"code"`
	Code   ErrorCode `json:"code"`
	Errors []string  `json:"errors"`
}

// メッセージがリッチ、ログは簡素
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

// サンプルstruct
type Internal struct {
}

func (i Internal) Internal() bool {
	return true
}

func (i Internal) Error() string {
	return "internal"
}

type ExternalError struct {
	ErrorMessage  string
	OriginalError error
	StatusCode    ErrorCode
}

func (e ExternalError) Messages() []string {
	return []string{e.Error()}
}

func (e ExternalError) Code() ErrorCode {
	return e.StatusCode
}

func (e ExternalError) Error() string {
	return e.OriginalError.Error()
}
