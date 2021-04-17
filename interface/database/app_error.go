package database

import (
	"net/http"
)

type ExternalError interface {
	Code() ErrorCode
	Messages() []string
}

type InternalError interface {
	Internal() bool
}

type ErrorCode int

const (
	AuthenticationParamMissing ErrorCode = iota
	AuthenticationFailure
	InvalidParameter
	UniqueValueDuplication
	ValueNotFound
	InternalServerError

	// Error codes for unhandled error
	UnHandledError ErrorCode = 999
)

var codeStatusMap = map[ErrorCode]int{
	AuthenticationFailure:      http.StatusForbidden,
	AuthenticationParamMissing: http.StatusBadRequest,
	InvalidParameter:           http.StatusBadRequest,
	UniqueValueDuplication:     http.StatusBadRequest,
	ValueNotFound:              http.StatusNotFound,
	InternalServerError:        http.StatusInternalServerError,
}

func GetHttpStatus(code ErrorCode) int {
	return codeStatusMap[code]
}

