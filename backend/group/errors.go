package group

import (
	"github.com/commonpool/backend/errors"
	"net/http"
)

var ErrInvalidGroupId = errors.NewWebServiceException("invalid group id", "ErrInvalidGroupId", http.StatusBadRequest)
var ErrMembershipPartyUnauthorized = errors.NewWebServiceException("not allowed to manage other people memberships", "ErrMembershipPartyUnauthorized", http.StatusForbidden)
var ErrManageMembershipsNotAdmin = errors.NewWebServiceException("don't have sufficient privilegtes", "ErrManageMembershipsNotAdmin", http.StatusForbidden)
