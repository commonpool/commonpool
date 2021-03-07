package test

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

var ApproveAllMatrix domain.PermissionMatrix = func(userKey keys.UserKey, offerItem domain.OfferItem, direction domain.ApprovalDirection) bool {
	return true
}

var DenyAllMatrix domain.PermissionMatrix = func(userKey keys.UserKey, offerItem domain.OfferItem, direction domain.ApprovalDirection) bool {
	return false
}
