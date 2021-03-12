package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type GetOfferPermissions struct {
	getOffersPermissions *GetOffersPermissions
}

func NewGetOfferPermissions(getOffersPermissions *GetOffersPermissions) *GetOfferPermissions {
	return &GetOfferPermissions{getOffersPermissions: getOffersPermissions}
}

func (q *GetOfferPermissions) getPermissionTuples(ctx context.Context, offerKey keys.OfferKey) (*PermissionTuples, error) {
	permissions, err := q.getOffersPermissions.Get(ctx, keys.NewOfferKeys([]keys.OfferKey{offerKey}))
	if err != nil {
		return nil, err
	}
	if permission, ok := permissions[offerKey]; ok {
		return permission, nil
	}
	return nil, exceptions.ErrOfferNotFound
}

func (q *GetOfferPermissions) Get(ctx context.Context, offerKey keys.OfferKey) (domain.OfferPermissionGetter, error) {
	return q.getPermissionTuples(ctx, offerKey)
}

type PermissionTuples struct {
	Permissions map[keys.OfferItemKey]*PermissionTuple
}

func NewPermissionTuples() *PermissionTuples {
	return &PermissionTuples{
		Permissions: map[keys.OfferItemKey]*PermissionTuple{},
	}
}

func (p *PermissionTuples) Can(userKey keys.UserKey, offerItem domain.OfferItem, direction domain.ApprovalDirection) bool {
	permissions, ok := p.Permissions[offerItem.GetKey()]
	if !ok {
		return false
	}
	if direction == domain.Inbound {
		return permissions.Inbound[userKey]
	}
	if direction == domain.Outbound {
		return permissions.Outbound[userKey]
	}
	return false
}

func (p *PermissionTuples) CanAny(userKey keys.UserKey) bool {
	for _, tuple := range p.Permissions {
		if tuple.Outbound[userKey] {
			return true
		}
		if tuple.Inbound[userKey] {
			return true
		}
	}
	return false
}

func (p *PermissionTuples) AddPermission(offerItemKey keys.OfferItemKey, userKey keys.UserKey, direction domain.ApprovalDirection) {
	if _, ok := p.Permissions[offerItemKey]; !ok {
		p.Permissions[offerItemKey] = NewPermissionTuple()
	}
	p.Permissions[offerItemKey].AddPermission(userKey, direction)
}

type PermissionTuple struct {
	Inbound  map[keys.UserKey]bool
	Outbound map[keys.UserKey]bool
}

func NewPermissionTuple() *PermissionTuple {
	return &PermissionTuple{
		Inbound:  map[keys.UserKey]bool{},
		Outbound: map[keys.UserKey]bool{},
	}
}

func (p *PermissionTuple) AddPermission(userKey keys.UserKey, direction domain.ApprovalDirection) {
	switch direction {
	case domain.Inbound:
		p.Inbound[userKey] = true
	case domain.Outbound:
		p.Outbound[userKey] = true
	}
}
