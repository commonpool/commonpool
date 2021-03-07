package keys

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type OfferPermissions struct {
	OfferKey    keys.OfferKey
	Permissions []OfferItemPermission
}

type OfferItemPermission struct {
	UserKey      keys.UserKey
	OfferItemKey keys.OfferItemKey
	Direction    domain.ApprovalDirection
}
