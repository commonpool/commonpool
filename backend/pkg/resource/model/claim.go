package model

type Claim struct {
	ResourceKey ResourceKey
	ClaimType   ClaimType
	For         *Target
}
