package errors

import "net/http"

type InternalError interface {
	Error() string
	Code() int
}

type internalError struct {
	code int
	err  string
}

func (a internalError) Error() string {
	return a.err
}

func (a internalError) Code() int {
	return a.code
}

func newInternalError(code int, err string) internalError {
	return internalError{
		code: code,
		err:  err,
	}
}

func BadRequestError(message string) InternalError {
	return newInternalError(http.StatusBadRequest, message)
}

func InternalServerErrorWithMsg(message string) InternalError {
	return newInternalError(http.StatusInternalServerError, message)
}

func ResourceNotFoundError(message string) InternalError {
	return newInternalError(http.StatusNotFound, message)
}

func ForbiddenError(message string) InternalError {
	return newInternalError(http.StatusForbidden, message)
}

var InternalServerError = newInternalError(http.StatusInternalServerError, "internal server error")
