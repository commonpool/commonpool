package resource

import "github.com/commonpool/backend/pkg/keys"

type Claim struct {
	ResourceKey keys.ResourceKey
	ClaimType   ClaimType
	For         *Target
}
