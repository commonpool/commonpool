package errors

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WebServiceException struct {
	Status  int
	Code    string
	Message string
}

func (e WebServiceException) Error() string {
	return e.Message
}

func (e WebServiceException) Is(err error) bool {
	a, ok := err.(WebServiceException)
	if !ok {
		return false
	}
	return e.Code == a.Code
}

func NewWebServiceException(message string, code string, status int) error {
	e := WebServiceException{
		Status:  status,
		Code:    code,
		Message: message,
	}
	return &e
}

type ErrorIs interface {
	Is(error) bool
}

var ErrUnauthorized = NewWebServiceException("unauthorized", "ErrUnauthorized", http.StatusUnauthorized)

func ReturnException(c echo.Context, err error) error {
	ws, ok := err.(*WebServiceException)
	if !ok {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Message:    "Internal server error",
			Code:       "ErrInternalServerError",
			StatusCode: http.StatusInternalServerError,
		})
	}
	return c.JSON(ws.Status, &ErrorResponse{
		Message:    ws.Message,
		Code:       ws.Code,
		StatusCode: ws.Status,
	})
}

var (
	ErrUserNotFound = func(user string) ErrorResponse {
		return NewError(fmt.Sprintf("user with id '%s' could not be found", user), "ErrUserNotFound", http.StatusNotFound)
	}
	ErrResourceNotFound = func(resource string) ErrorResponse {
		return NewError(fmt.Sprintf("resource with id '%s' could not be found", resource), "ErrResourceNotFound", http.StatusNotFound)
	}
	ErrCreateResourceBadRequest = func(err error) ErrorResponse {
		return NewError("could not process create resource request: "+err.Error(), "ErrCreateResourceBadRequest", http.StatusBadRequest)
	}
	ErrUpdateResourceBadRequest = func(err error) ErrorResponse {
		return NewError("could not process update resource request: "+err.Error(), "ErrUpdateResourceBadRequest", http.StatusBadRequest)
	}
	ErrSendResourceMsgBadRequest = func(err error) ErrorResponse {
		return NewError("could not process send message request: "+err.Error(), "ErrSendResourceMsgBadRequest", http.StatusBadRequest)
	}
	ErrSendOfferBadRequest = func(err error) ErrorResponse {
		return NewError("could not process send offer request: "+err.Error(), "ErrSendOfferBadRequest", http.StatusBadRequest)
	}
	ErrValidation = func(msg string) *ErrorResponse {
		err := NewError("validation error: "+msg, "ErrValidation", http.StatusBadRequest)
		return &err
	}
	ErrInvalidResourceKey = func(key string) ErrorResponse {
		return NewError(fmt.Sprintf("invalid resource key: '%s'", key), "ErrInvalidResourceKey", http.StatusBadRequest)
	}
	ErrParseSkip = func(err string) ErrorResponse {
		return NewError(fmt.Sprintf("cannot parse skip: '%s'", err), "ErrParseSkip", http.StatusBadRequest)
	}
	ErrParseTake = func(err string) ErrorResponse {
		return NewError(fmt.Sprintf("cannot parse take: '%s'", err), "ErrParseTake", http.StatusBadRequest)
	}
	ErrParseBefore = func(err string) ErrorResponse {
		return NewError(fmt.Sprintf("cannot parse before: '%s'", err), "ErrParseBefore", http.StatusBadRequest)
	}
	ErrParseResourceType = func(resType string) ErrorResponse {
		return NewError(fmt.Sprintf("cannot parse resource type '%s'", resType), "ErrParseResourceType", http.StatusBadRequest)
	}
	ErrCannotConvertToInt = func(int string, err string) ErrorResponse {
		return NewError(fmt.Sprintf("cannot convert '%s' to integer: %s", int, err), "ErrCannotParseInt", http.StatusBadRequest)
	}
	ErrCannotInquireAboutOwnResource = func() *ErrorResponse {
		err := NewError("cannot inquire about your own resource", "ErrCannotInquireAboutOwnResource", http.StatusForbidden)
		return &err
	}
	ErrInvalidTopicId = func(threadId string) ErrorResponse {
		return NewError(fmt.Sprintf("invalid thread id: '%s'", threadId), "ErrInvalidTopicId", http.StatusBadRequest)
	}
	ErrTransactionResourceOwnerMismatch = func() ErrorResponse {
		return NewError("resource owner doesn't match offer.item[].from", "ErrTransactionResourceOwnerMismatch", http.StatusBadRequest)
	}
)

func NewError(message string, code string, statusCode int) ErrorResponse {
	return ErrorResponse{
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
	}
}

type ErrorResponse struct {
	Message    string
	Code       string
	StatusCode int
}

func (r *ErrorResponse) Error() string {
	return r.Message
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
