package errorX

import (
	"errors"
	"net/http"
)

// ErrorCode a list of failure codes, corresponding to the
// database field "failure_code" and the API response
// "failure_code"
type ErrorCode int

const (
	ERROR_CODE_NOT_AUTHENTICATED ErrorCode = iota + 4000
	ERROR_CODE_REFERAL_NOT_FOUND_OR_EXPIRED
	ERROR_CODE_FORBIDDEN_GENERATE_LINK
)
const (
	ERROR_CODE_INTERNAL_SERVER ErrorCode = iota + 5000
)

var ERROR_HTTP_MAP = map[ErrorCode]int{
	ERROR_CODE_INTERNAL_SERVER:              http.StatusInternalServerError,
	ERROR_CODE_NOT_AUTHENTICATED:            http.StatusUnauthorized,
	ERROR_CODE_REFERAL_NOT_FOUND_OR_EXPIRED: http.StatusNotFound,
	ERROR_CODE_FORBIDDEN_GENERATE_LINK:      http.StatusForbidden,
}

var ERROR_MAP = map[ErrorCode]error{
	ERROR_CODE_INTERNAL_SERVER:              errors.New("Something went wrong please contact developer"),
	ERROR_CODE_NOT_AUTHENTICATED:            errors.New(http.StatusText(http.StatusUnauthorized)),
	ERROR_CODE_REFERAL_NOT_FOUND_OR_EXPIRED: errors.New("Code not found or expired"),
	ERROR_CODE_FORBIDDEN_GENERATE_LINK:      errors.New("only user with role generator that can generate link"),
}
