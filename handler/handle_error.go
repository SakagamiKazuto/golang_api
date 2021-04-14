package handler

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/apperror"
	"github.com/labstack/echo/v4"
)

type ExternalHandleError struct {
	ErrorMessage  string
	OriginalError error
	StatusCode    apperror.ErrorCode
}

func (e ExternalHandleError) Messages() []string {
	return []string{e.Error()}
}

func (e ExternalHandleError) Code() apperror.ErrorCode {
	return e.StatusCode
}

func (e ExternalHandleError) Error() string {
	return e.ErrorMessage + "\n" + e.OriginalError.Error()
}

func createLoginFailureErr(c echo.Context, err error) error {
	return apperror.ResponseError(c, ExternalHandleError{
		ErrorMessage:  fmt.Sprintf("ログインに失敗しました"),
		OriginalError: err,
		StatusCode:    apperror.AuthenticationFailure,
	})
}
