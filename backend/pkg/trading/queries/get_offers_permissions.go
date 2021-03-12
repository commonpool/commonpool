package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	groupreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strings"
)

type GetOffersPermissions struct {
	db                 *gorm.DB
	getOffersGroupKeys *GetOffersGroupKeys
}

func NewGetOffersPermissions(db *gorm.DB, getOffersGroupKeys *GetOffersGroupKeys) *GetOffersPermissions {
	return &GetOffersPermissions{db: db, getOffersGroupKeys: getOffersGroupKeys}
}

type resourceResult struct {
	ResourceKey keys.ResourceKey
	Owner       keys.Target `json:"owner" gorm:"embedded;embeddedPrefix:owner_"`
}

func (q *GetOffersPermissions) getPermissionTuples(ctx context.Context, offerKeys *keys.OfferKeys) (map[keys.OfferKey]*PermissionTuples, error) {

	if len(offerKeys.Items) == 0 {
		return map[keys.OfferKey]*PermissionTuples{}, nil
	}

	var offersSb strings.Builder
	var offersParams []interface{}
	offersSb.WriteString("offer_key in (")
	for i, offerKey := range offerKeys.Items {
		offersSb.WriteString("?")
		offersParams = append(offersParams, offerKey.String())
		if i < len(offerKeys.Items)-1 {
			offersSb.WriteString(",")
		}
	}

	offersSb.WriteString(")")

	var offerGroups map[keys.OfferKey]keys.GroupKey
	var offerItems []*groupreadmodels.OfferItemReadModel
	var resources []*resourceResult

	// Find offers and offeritems for offer keys
	g1, ctx := errgroup.WithContext(ctx)
	g1.Go(func() error {
		var err error
		offerGroups, err = q.getOffersGroupKeys.Get(ctx, offerKeys)
		return err
	})
	g1.Go(func() error {
		return q.db.Model(&groupreadmodels.OfferItemReadModel{}).Where(offersSb.String(), offersParams...).Find(&offerItems).Find(&offerItems).Error
	})
	g1.Go(func() error {
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
			where `+offersSb.String(), offersParams...).
			Scan(&resources).
			Error
	})

	if err := g1.Wait(); err != nil {
		return nil, err
	}

	if len(offerGroups) == 0 {
		return map[keys.OfferKey]*PermissionTuples{}, nil
	}

	groupAdmins, err := q.getAdminsForOffers(ctx, offerGroups)
	if err != nil {
		return nil, err
	}

	adminsByGroup := map[keys.GroupKey][]*groupreadmodels.OfferUserMembershipReadModel{}
	for _, groupAdmin := range groupAdmins {
		adminsByGroup[groupAdmin.GroupKey] = append(adminsByGroup[groupAdmin.GroupKey], groupAdmin)
	}

	var groupedOfferItems = map[keys.OfferKey][]*groupreadmodels.OfferItemReadModel{}
	for _, offerItem := range offerItems {
		groupedOfferItems[offerItem.OfferKey] = append(groupedOfferItems[offerItem.OfferKey], offerItem)
	}

	var resourceMap = map[keys.ResourceKey]*resourceResult{}
	for _, resource := range resources {
		resourceMap[resource.ResourceKey] = resource
	}

	var permissionMap = map[keys.OfferKey]*PermissionTuples{}

	for offerKey, groupKey := range offerGroups {
		var permissions = NewPermissionTuples()
		offerItems := groupedOfferItems[offerKey]
		adminsForOfferGroup := adminsByGroup[groupKey]
		for _, offerItem := range offerItems {
			if offerItem.From.IsUser() {
				permissions.AddPermission(offerItem.OfferItemKey, offerItem.From.GetUserKey(), domain.Outbound)
			}
			if offerItem.To.IsUser() {
				permissions.AddPermission(offerItem.OfferItemKey, offerItem.To.GetUserKey(), domain.Inbound)
			}
			if offerItem.From.IsGroup() {
				for _, admin := range adminsForOfferGroup {
					permissions.AddPermission(offerItem.OfferItemKey, admin.UserKey, domain.Outbound)
				}
			}
			if offerItem.To.IsGroup() {
				for _, admin := range adminsForOfferGroup {
					permissions.AddPermission(offerItem.OfferItemKey, admin.UserKey, domain.Inbound)
				}
			}
			if offerItem.ResourceKey != nil {
				if resource, ok := resourceMap[*offerItem.ResourceKey]; ok {
					if resource.Owner.IsUser() {
						permissions.AddPermission(offerItem.OfferItemKey, resource.Owner.GetUserKey(), domain.Outbound)
					} else if resource.Owner.IsGroup() {
						for _, adminForOfferGroup := range adminsForOfferGroup {
							if resource.Owner.GetGroupKey() == adminForOfferGroup.GroupKey {
								permissions.AddPermission(offerItem.OfferItemKey, adminForOfferGroup.UserKey, domain.Outbound)
							}
						}
					}

				}
			}
		}
		permissionMap[offerKey] = permissions
	}

	return permissionMap, nil
}

func (q *GetOffersPermissions) getAdminsForOffers(ctx context.Context, offers map[keys.OfferKey]keys.GroupKey) ([]*groupreadmodels.OfferUserMembershipReadModel, error) {

	// Get unique group keys for offers
	var groups = map[keys.GroupKey]bool{}
	for _, groupKey := range offers {
		groups[groupKey] = true
	}

	var groupKeySet []string
	for groupKey, _ := range groups {
		groupKeySet = append(groupKeySet, groupKey.String())
	}

	// build sql for querying admins
	var adminSb strings.Builder
	adminSb.WriteString("group_key in (")
	var adminParams []interface{}
	for i, groupKey := range groupKeySet {
		adminSb.WriteString("?")
		if i < len(groupKeySet)-1 {
			adminSb.WriteString(",")
		}
		adminParams = append(adminParams, groupKey)
	}
	adminSb.WriteString(") and (is_admin = ? or is_owner = ?)")
	adminParams = append(adminParams, true, true)

	var groupAdmins []*groupreadmodels.OfferUserMembershipReadModel

	if err := q.db.Model(&groupreadmodels.OfferUserMembershipReadModel{}).
		Where(adminSb.String(), adminParams...).
		Find(&groupAdmins).
		Error; err != nil {
		return nil, err
	}

	return groupAdmins, nil
}

func (q *GetOffersPermissions) Get(ctx context.Context, offerKeys *keys.OfferKeys) (map[keys.OfferKey]*PermissionTuples, error) {
	return q.getPermissionTuples(ctx, offerKeys)
}
