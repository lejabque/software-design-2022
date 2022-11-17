package controllers

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	message  string
	httpcode int
}

func (e *CustomError) Error() string {
	return e.message
}

func (e *CustomError) Code() int {
	return e.httpcode
}

func NewCustomError(httpcode int, format string, params ...any) *CustomError {
	return &CustomError{message: fmt.Sprintf(format, params...), httpcode: httpcode}
}

func NewUserError(format string, params ...any) *CustomError {
	return NewCustomError(http.StatusBadRequest, format, params...)
}

func ErrorCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Code()
	}
	return http.StatusInternalServerError
}
