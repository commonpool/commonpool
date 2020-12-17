package resource

import (
	"github.com/commonpool/backend/pkg/keys"
)

type Claims struct {
	Items []*Claim
}

func NewClaims(items []*Claim) *Claims {
	copied := make([]*Claim, len(items))
	copy(copied, items)
	return &Claims{
		Items: copied,
	}
}

func NewEmptyClaims() *Claims {
	return &Claims{
		Items: []*Claim{},
	}
}

func (c *Claims) AppendAll(claims *Claims) {
	for _, claim := range claims.Items {
		c.Items = append(c.Items, claim)
	}
}

func (c *Claims) UserHasClaim(userKey keys.UserKey, resourceKey keys.ResourceKey, claimType ClaimType) bool {
	for _, claim := range c.Items {
		if claim.ClaimType == claimType && claim.ResourceKey == resourceKey && claim.For.IsForUser() && claim.For.GetUserKey() == userKey {
			return true
		}
	}
	return false
}

func (c *Claims) GroupHasClaim(groupKey keys.GroupKey, resourceKey keys.ResourceKey, claimType ClaimType) bool {
	for _, claim := range c.Items {
		if claim.ClaimType == claimType && claim.ResourceKey == resourceKey && claim.For.IsForGroup() && claim.For.GetGroupKey() == groupKey {
			return true
		}
	}
	return false
}

func (c *Claims) HasClaim(target *Target, resourceKey keys.ResourceKey, claimType ClaimType) bool {
	for _, claim := range c.Items {
		if claim.ClaimType == claimType &&
			claim.ResourceKey == resourceKey &&
			claim.For.Equals(target) {
			return true
		}
	}
	return false
}
