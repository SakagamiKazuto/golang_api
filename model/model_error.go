package model

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/apperror"
	"path/filepath"
	"runtime"
)

type InternalDBError struct {
	ErrorMessage  string
	OriginalError error
}

func (i InternalDBError) Internal() bool {
	return true
}

func (i InternalDBError) Error() string {
	return i.ErrorMessage + i.OriginalError.Error()
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

func createInErrMsg(skip int) string {
	pc, file, _, ok := runtime.Caller(skip)
	if ok {
		fname := filepath.Base(file)
		return fmt.Sprintf(`DB処理中にエラーが発生しました\nfile: %s, func: %s\n`,
			fname, runtime.FuncForPC(pc).Name())
	}
	return fmt.Sprintf(`DB処理中にエラーが発生しました\n`)
}