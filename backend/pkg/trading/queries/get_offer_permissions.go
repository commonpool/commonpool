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

func (q *GetOfferPermissions) getPermissionTuples(ctx context.Context, offerKey keys.OfferKey) (*PermissionTuples, error) {

	g, ctx := errgroup.WithContext(ctx)

	var offer groupreadmodels.DBOfferReadModel
	if err := q.db.Model(&groupreadmodels.DBOfferReadModel{}).Where("offer_key = ?", offerKey).First(&offer).Error; err != nil {
		return nil, err
	}

	var offerItems []*groupreadmodels.OfferItemReadModel
	var groupAdmins []*groupreadmodels.OfferUserMembershipReadModel

	g.Go(func() error {
		return q.db.Model(&groupreadmodels.OfferItemReadModel{}).Find(&offerItems, "offer_key = ?", offerKey).Find(&offerItems).Error
	})

	g.Go(func() error {
		return q.db.Model(&groupreadmodels.OfferUserMembershipReadModel{}).
			Where("group_key = ? and (is_member = ? or is_owner = ?)", offer.GroupKey, true, true).
			Find(&groupAdmins).
			Error
	})

	type resourceResult struct {
		ResourceKey keys.ResourceKey
		Owner       keys.Target `json:"owner" gorm:"embedded;embeddedPrefix:owner_"`
	}
	var resources []*resourceResult
	g.Go(func() error {
		return q.db.Raw(`
			select
				offer_resource_read_models.resource_key,
				offer_resource_read_models.owner_user_key,
				offer_resource_read_models.owner_group_key,
				offer_resource_read_models.owner_type
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
			permissions.AddPermission(offerItem.OfferItemKey, offerItem.To.GetUserKey(), domain.Inbound)
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
