package database

import (
	"fmt"
	"github.com/lib/pq"
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
	return fmt.Sprintf("Message: %s\nDetail: %s\nPlace: %s %s", i.Message, i.Detail, i.File, i.Line)
}

type ExternalDBError struct {
	ErrorMessage  string
	OriginalError error
	StatusCode    ErrorCode
}

func (e ExternalDBError) Messages() []string {
	return []string{e.Error()}
}

func (e ExternalDBError) Code() ErrorCode {
	return e.StatusCode
}

func (e ExternalDBError) Error() string {
	return e.ErrorMessage + "\n" + e.OriginalError.Error()
}

func CreateInDBError(err error) error {
	pqe, ok := err.(*pq.Error)

	if !ok {
		return err
	}

	return &InternalDBError{
		Message:       pqe.Message,
		Detail:        pqe.Detail,
		File:          pqe.File,
		Line:          pqe.Line,
		OriginalError: pqe,
	}
}

