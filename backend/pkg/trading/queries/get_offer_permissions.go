package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	groupreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type GetOfferPermissions struct {
	db *gorm.DB
}

func NewGetOfferPermissions(db *gorm.DB) *GetOfferPermissions {
	return &GetOfferPermissions{db: db}
}

func (q *GetOfferPermissions) Get(ctx context.Context, offerKey keys.OfferKey) (domain.OfferPermissionGetter, error) {

	g, ctx := errgroup.WithContext(ctx)

	var offerItems []*groupreadmodels.OfferItemReadModel
	var groupAdmins []*groupreadmodels.OfferUserMembershipReadModel

	g.Go(func() error {
		return q.db.Model(&offerItems).Find("offer_key = ?", offerKey).Find(&offerItems).Error
	})

	g.Go(func() error {
		return q.db.Model(&groupreadmodels.OfferUserMembershipReadModel{}).
			Where("user_key = ? and group_key = ? and (is_member = true or is_owner = true)").
			Find(&groupAdmins).
			Error
	})

	type resourceResult struct {
		ResourceKey keys.ResourceKey
		Owner       domain.Target
	}
	var resources []*resourceResult
	g.Go(func() error {
		return q.db.Raw(`
			select
				offer_resource_read_models.resource_key,
				offer_resource_read_models.owner_user_key,
				offer_resource_read_models.owner_group_key
			from 
				offer_item_read_models
			join 
				offer_resource_read_models on offer_resource_read_models.resource_key = offer_item_read_models.resource_key
			where 
				offer_key = ?`,
			offerKey).
			Scan(&resources).
			Error
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var permissions = NewPermissionTuples()

	for _, offerItem := range offerItems {
		if offerItem.From.IsForUser() {
			permissions.AddPermission(offerItem.OfferItemKey, offerItem.From.GetUserKey(), domain.Outbound)
		}
		if offerItem.To.IsForUser() {
			permissions.AddPermission(offerItem.OfferItemKey, offerItem.To.GetUserKey(), domain.Outbound)
		}
		if offerItem.From.IsForGroup() {
			for _, groupAdmin := range groupAdmins {
				permissions.AddPermission(offerItem.OfferItemKey, groupAdmin.UserKey, domain.Outbound)
			}
		}
		if offerItem.To.IsForGroup() {
			for _, groupAdmin := range groupAdmins {
				permissions.AddPermission(offerItem.OfferItemKey, groupAdmin.UserKey, domain.Inbound)
			}
		}
		if offerItem.ResourceKey != nil {
			// find resource
			for _, resource := range resources {
				if resource.ResourceKey == *offerItem.ResourceKey {
					if resource.Owner.IsForUser() {
						permissions.AddPermission(offerItem.OfferItemKey, resource.Owner.GetUserKey(), domain.Outbound)
					} else if resource.Owner.IsForGroup() {
						for _, groupAdmin := range groupAdmins {
							if resource.Owner.GetGroupKey() == groupAdmin.GroupKey {
								permissions.AddPermission(offerItem.OfferItemKey, groupAdmin.UserKey, domain.Outbound)
							}
						}
					}
				}
			}
		}
	}

	return permissions, nil

}

type PermissionTuples struct {
	Permissions map[keys.OfferItemKey]*PermissionTuple
}

func NewPermissionTuples() *PermissionTuples {
	return &PermissionTuples{
		Permissions: map[keys.OfferItemKey]*PermissionTuple{},
	}
}

func (p *PermissionTuples) Can(userKey keys.UserKey, key keys.OfferItemKey, direction domain.ApprovalDirection) bool {
	permissions, ok := p.Permissions[key]
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

func (p *PermissionTuples) AddPermission(offerItemKey keys.OfferItemKey, userKey keys.UserKey, direction domain.ApprovalDirection) {
	if _, ok := p.Permissions[offerItemKey]; !ok {
		p.Permissions[offerItemKey] = NewPermissionTuple()
	}
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
