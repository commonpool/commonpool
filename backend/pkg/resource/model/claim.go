package model

import "github.com/commonpool/backend/model"

type Claim struct {
	ResourceKey ResourceKey
	ClaimType   ClaimType
	For         *model.Target
}
