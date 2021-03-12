package domain

import (
	"github.com/commonpool/backend/pkg/keys"
)

type OfferItem interface {
	GetType() OfferItemType
	GetOfferKey() keys.OfferKey
	GetKey() keys.OfferItemKey
}

type ResourceKeyGetter interface {
	GetResourceKey() keys.ResourceKey
}

type FromTargeter interface {
	GetFrom() keys.Target
}

type ToTargeter interface {
	GetTo() keys.Target
}

type ResourceOfferItem interface {
	ToTargeter
	ResourceKeyGetter
}
