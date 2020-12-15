package exceptions

import (
	"fmt"
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

var ErrUserNotFound = NewWebServiceException("user not found", "ErrUserNotFound", http.StatusNotFound)
var ErrResourceNotFound = NewWebServiceException("resource not found", "ErrResourceNotFound", http.StatusNotFound)
var ErrGroupNotFound = NewWebServiceException("group not found", "ErrGroupNotFound", http.StatusNotFound)
var ErrOfferNotFound = NewWebServiceException("offer not found", "ErrOfferNotFound", http.StatusNotFound)
var ErrOfferItemNotFound = NewWebServiceException("offer item not found", "ErrOfferItemNotFound", http.StatusNotFound)
var ErrMembershipNotFound = NewWebServiceException("membership not found", "ErrMembershipNotFound", http.StatusNotFound)
var ErrUserOrGroupNotFound = NewWebServiceException("user or group not found", "ErrUserOrGroupNotFound", http.StatusNotFound)
var ErrUnknownParty = NewWebServiceException("unknown party", "ErrUnknownParty", http.StatusBadRequest)
var ErrNegativeDuration = NewWebServiceException("time offers must have positive time value", "ErrNegativeDuration", http.StatusBadRequest)
var ErrWrongOfferItemType = NewWebServiceException("wrong offer item type", "ErrWrongOfferItemType", http.StatusBadRequest)
var ErrUnauthorized = NewWebServiceException("unauthorized", "ErrUnauthorized", http.StatusUnauthorized)
var ErrForbidden = NewWebServiceException("forbidden", "ErrForbidden", http.StatusForbidden)
var ErrDuplicateResourceInOffer = NewWebServiceException("resource can only appear once in an offer", "ErrDuplicateResourceInOffer", http.StatusBadRequest)
var ErrResourceNotSharedWithGroup = NewWebServiceException("resource is not shared with the group", "ErrResourceNotSharedWithGroup", http.StatusBadRequest)
var ErrCannotTransferResourceToItsOwner = NewWebServiceException("resource cannot be transferred to its own owner", "ErrCannotTransferResourceToItsOwner", http.StatusBadRequest)
var ErrResourceTransferOfferItemsMustReferToObjectResources = NewWebServiceException("resource transfers can only be for object-typed resources", "ErrResourceTransferOfferItemsMustReferToObjectResources", http.StatusBadRequest)
var ErrServiceProvisionOfferItemsMustPointToServiceResources = NewWebServiceException("service provision offer items must be for a service-type resource!", "ErrServiceProvisionOfferItemsMustPointToServiceResources", http.StatusBadRequest)
var ErrBorrowOfferItemMustReferToObjectTypedResource = NewWebServiceException("borrow offer items must be for a service-type resource!", "ErrBorrowOfferItemMustReferToObjectTypedResource", http.StatusBadRequest)
var ErrInvalidOfferItemType = NewWebServiceException("invalid offer item type", "ErrInvalidOfferItemType", http.StatusBadRequest)
var ErrInvalidTargetType = NewWebServiceException("invalid target type", "ErrInvalidTargetType", http.StatusBadRequest)
var ErrInvalidTakeQueryParam = NewWebServiceException("query parameter 'take' is invalid", "ErrInvalidTakeQueryParam", http.StatusBadRequest)
var ErrInvalidSkipQueryParam = NewWebServiceException("query parameter 'skip' is invalid", "ErrInvalidSkipQueryParam", http.StatusBadRequest)
var ErrInvalidBeforeQueryParam = NewWebServiceException("query parameter 'before' is invalid", "ErrInvalidBeforeQueryParam", http.StatusBadRequest)
var ErrInvalidGroupId = NewWebServiceException("invalid group id", "ErrInvalidGroupId", http.StatusBadRequest)
var ErrMembershipPartyUnauthorized = NewWebServiceException("not allowed to manage other people memberships", "ErrMembershipPartyUnauthorized", http.StatusForbidden)
var ErrManageMembershipsNotAdmin = NewWebServiceException("don't have sufficient privilegtes", "ErrManageMembershipsNotAdmin", http.StatusForbidden)

func ErrQueryParamRequired(queryParameter string) error {
	return NewWebServiceException(fmt.Sprintf("query parameter '%s' is required", queryParameter), "ErrQueryParamRequired", http.StatusBadRequest)
}

type ValidError struct {
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Namespace       string `json:"namespace"`
	StructNamespace string `json:"structNamespace"`
	Field           string `json:"field"`
	StructField     string `json:"structField"`
	Param           string `json:"param"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
}

type ValidErrors struct {
	Errors  []ValidError      `json:"errors"`
	Message string            `json:"message"`
	Trans   map[string]string `json:"trans"`
}

var (
	ErrValidation = func(msg string) *ErrorResponse {
		err := NewError("validation error: "+msg, "ErrValidation", http.StatusBadRequest)
		return &err
	}
	ErrInvalidResourceKey = func(key string) ErrorResponse {
		return NewError(fmt.Sprintf("invalid resource key: '%s'", key), "ErrInvalidResourceKey", http.StatusBadRequest)
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
	Validation map[string]string
	Errors     []ValidError
}

func (r *ErrorResponse) Error() string {
	return r.Message
}
func (r *ErrorResponse) IsNotFoundError() bool {
	return r.StatusCode == http.StatusNotFound
}