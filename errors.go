package app

import (
	"errors"
	"fmt"
	"net/http"
)

// Application errors.
var (
	// General
	ErrNotFound              = errors.New("Not found")
	ErrInternal              = errors.New("Internal")
	ErrBadRequest            = errors.New("Bad request")
	ErrValidation            = errors.New("Invalid request body")
	ErrInvalidParamType      = errors.New("Invalid param type")
	ErrNotAllRequiredQueries = errors.New("Not all queries")
)

var errorCodesMap = map[error]int{
	ErrNotFound:              404,
	ErrInternal:              500,
	ErrBadRequest:            400,
	ErrValidation:            3,
	ErrNotAllRequiredQueries: 5,
}

var codesToErrorsMap = map[int]error{
	404: ErrNotFound,
	500: ErrInternal,
	400: ErrBadRequest,
	3:   ErrValidation,
	5:   ErrNotAllRequiredQueries,
}

func WrapE(err error, msg string) error {
	return fmt.Errorf("%w. Info: "+msg, err)
}

func GetHTTPCodeFromError(err error) int {
	switch err {
	case ErrInternal:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadRequest, ErrValidation, ErrInvalidParamType:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}

func ErrorCode(err error) int {
	for k, v := range errorCodesMap {
		if errors.Is(err, k) {
			return v
		}
	}
	return 1
}

func CodeToError(code int) (error, bool) {
	err, ok := codesToErrorsMap[code]
	return err, ok
}
