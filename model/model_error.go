package model

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/apperror"
)

type InternalDBError struct {
	Message       string
	Detail        string
	File        string
	Line        string
	OriginalError error
}

func (i InternalDBError) Internal() bool {
	return true
}

func (i InternalDBError) Error() string {
	return fmt.Sprintf(`Message: %s
Detail: %s
Place: %s %s`, i.Message, i.Detail, i.File, i.Line)
}

type ExternalDBError struct {
	ErrorMessage  string
	OriginalError error
	StatusCode    apperror.ErrorCode
}

func (e ExternalDBError) Messages() []string {
	return []string{e.Error()}
}

func (e ExternalDBError) Code() apperror.ErrorCode {
	return e.StatusCode
}

func (e ExternalDBError) Error() string {
	return e.ErrorMessage + ":" + e.OriginalError.Error()
}
