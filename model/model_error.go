package model

import "github.com/SakagamiKazuto/golang_api/apperror"

type InternalDBError struct {
}

func (i InternalDBError) Internal() bool {
	return true
}

func (i InternalDBError) Error() string {
	return "internal"
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
	return e.OriginalError.Error()
}
