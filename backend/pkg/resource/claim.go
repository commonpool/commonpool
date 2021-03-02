package resource

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

type Claim struct {
	ResourceKey keys.ResourceKey
	ClaimType   ClaimType
	For         *trading.Target
}
