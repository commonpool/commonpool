package resource

import "github.com/commonpool/backend/model"

type Claim struct {
	ResourceKey model.ResourceKey
	ClaimType   ClaimType
	For         *model.Target
}
