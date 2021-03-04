package domain

import "github.com/commonpool/backend/pkg/keys"

type PermissionMatrix func(userKey keys.UserKey, offerItem OfferItem, direction ApprovalDirection) bool

type RolePermissionMatrix struct {
	groupOwners    keys.UserKey
	groupAdmins    keys.UserKeys
	resourceOwners map[keys.ResourceKey]keys.UserKey
}
