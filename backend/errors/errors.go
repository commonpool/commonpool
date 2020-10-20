package errors

import (
	"fmt"
	"net/http"
)

const (
	ErrUserNotFoundMsg                 = "User with id %s could not be found"
	ErrUserNotFoundCode                = "USR-1000"
	ErrResourceNotFoundMsg             = "Resource with id %s could not be found"
	ErrResourceNotFoundCode            = "RES-PER-1000"
	ErrCreateResourceCannotBind        = "Cannot process resource"
	ErrCreateResourceCannotBindCode    = "RES-WEB-1000"
	ErrUpdateResourceCannotBind        = "Cannot process resource"
	ErrUpdateResourceCannotBindCode    = "RES-WEB-1001"
	ErrSummaryEmptyOrNull              = "Summary cannot be empty or null"
	ErrSummaryEmptyOrNullCode          = "WEB-VAL-1000"
	ErrSummaryTooLong                  = "Summary is too long"
	ErrSummaryTooLongCode              = "WEB-VAL-1001"
	ErrDescriptionTooLong              = "Description is too long"
	ErrDescriptionTooLongCode          = "WEB-VAL-1002"
	ErrUuidParseError                  = "Could not parse uuid"
	ErrUuidParseErrorCode              = "WEB-VAL-1003"
	ErrInternalServerMsg               = "Internal server error"
	ErrInternalServerCode              = "S5000"
	ErrExchangeValueTooLow             = "Exchange value too low"
	ErrExchangeValueTooLowCode         = "WEB-VAL-1004"
	ErrExchangeValueTooHigh            = "Exchange value too high"
	ErrExchangeValueTooHighCode        = "WEB-VAL-1005"
	ErrTimeSensitivityValueTooLow      = "TimeSensitivity value too low"
	ErrTimeSensitivityValueTooLowCode  = "WEB-VAL-1006"
	ErrTimeSensitivityValueTooHigh     = "TimeSensitivity value too high"
	ErrTimeSensitivityValueTooHighCode = "WEB-VAL-1007"
	ErrNecessityLevelValueTooLow       = "NecessityLevel value too low"
	ErrNecessityLevelValueTooLowCode   = "WEB-VAL-1008"
	ErrNecessityLevelValueTooHigh      = "NecessityLevel value too high"
	ErrNecessityLevelValueTooHighCode  = "WEB-VAL-1009"
	ErrInvalidResourceType             = "Invalid resource type"
	ErrInvalidResourceTypeCode         = "WEB-VAL-1010"
	ErrInvalidTake                     = "'take' query parameter cannot be converted to int"
	ErrInvalidTakeCode                 = "WEB-VAL-1011"
	ErrInvalidSkip                     = "'skip' query parameter cannot be converted to int"
	ErrInvalidSkipCode                 = "WEB-VAL-1012"
)

type ErrorResponse struct {
	Message    string
	Code       string
	StatusCode int
}

func (r *ErrorResponse) Error() string {
	return r.Message
}

func NewError(message string, code string, statusCode int) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
	}
}

func NewResourceNotFoundError(key string) *ErrorResponse {
	return NewError(
		fmt.Sprintf(ErrResourceNotFoundMsg, key),
		ErrResourceNotFoundCode,
		http.StatusNotFound)
}

func NewUserNotFoundError(key string) *ErrorResponse {
	return NewError(
		fmt.Sprintf(ErrUserNotFoundMsg, key),
		ErrUserNotFoundCode,
		http.StatusNotFound)
}

func NewInternalServerError() *ErrorResponse {
	return NewError(
		ErrInternalServerMsg,
		ErrInternalServerCode,
		http.StatusInternalServerError)
}

func (r *ErrorResponse) IsNotFoundError() bool {
	return r.StatusCode == http.StatusNotFound
}

func IsNotFoundError(err error) bool {
	res, ok := err.(*ErrorResponse)
	if !ok {
		return false
	}
	return res.StatusCode == http.StatusNotFound
}
