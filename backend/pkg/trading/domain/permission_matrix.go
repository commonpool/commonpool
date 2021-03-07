package domain

import (
	"github.com/commonpool/backend/pkg/keys"
)

type OfferPermissionGetter interface {
	Can(userKey keys.UserKey, offerItem OfferItem, direction ApprovalDirection) bool
}

type PermissionMatrix func(userKey keys.UserKey, offerItem OfferItem, direction ApprovalDirection) bool

type RolePermissionMatrix struct {
	groupOwners    keys.UserKey
	groupAdmins    keys.UserKeys
	resourceOwners map[keys.ResourceKey]keys.UserKey
}