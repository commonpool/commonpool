package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
	"strings"
)

type GetOfferKeysForUser struct {
	db                           *gorm.DB
	getUserAdministeredGroupKeys *GetUserAdministeredGroupKeys
}

func NewGetOfferKeysForUser(db *gorm.DB, getUserAdministeredGroupKeys *GetUserAdministeredGroupKeys) *GetOfferKeysForUser {
	return &GetOfferKeysForUser{
		db:                           db,
		getUserAdministeredGroupKeys: getUserAdministeredGroupKeys,
	}
}

func (q *GetOfferKeysForUser) Get(ctx context.Context, userKey keys.UserKey) (*keys.OfferKeys, error) {

	// find groups where user is admin/owner

	administeredGroupKeys, err := q.getUserAdministeredGroupKeys.Get(ctx, userKey)
	if err != nil {
		return nil, err
	}

	/**
	select oi.offer_key
	from offer_item_read_models oi
	left join offer_read_models o on o.offer_key = oi.offer_key
	left join offer_resource_read_models r on oi.resource_key = r.resource_key
	where
	    oi.from_user_key = 'b0819e0d-fdf8-4c90-90d5-f359d2c079a4'
		or oi.to_user_key = 'b0819e0d-fdf8-4c90-90d5-f359d2c079a4'
		or from_group_key in ('90448a69-9bf3-4c97-8287-e8843ccb3d7f')

	   or r.owner_user_key = 'b0819e0d-fdf8-4c90-90d5-f359d2c079a4'
	   or from_group_key in ('90448a69-9bf3-4c97-8287-e8843ccb3d7f')
	   or oi.to_group_key in ('90448a69-9bf3-4c97-8287-e8843ccb3d7f')
	   or r.owner_group_key in ('90448a69-9bf3-4c97-8287-e8843ccb3d7f')
	*/

	var groupsSb strings.Builder
	var administeredGroupCount = administeredGroupKeys.Count()
	var params = make([]interface{}, 3+3*administeredGroupCount)
	params[0] = userKey
	params[1] = userKey
	params[2] = userKey

	if administeredGroupCount > 0 {
		groupsSb.WriteString("(")
		for i := 0; i < administeredGroupCount; i++ {
			for j := 0; j < 3; j++ {
				params[j*administeredGroupCount+i+3] = administeredGroupKeys.Items[i].String()
			}
			groupsSb.WriteString("?")
			if i < administeredGroupCount-1 {
				groupsSb.WriteString(",")
			}
		}
		groupsSb.WriteString(")")
	}

	var sb strings.Builder

	sb.WriteString("select oi.offer_key ")
	sb.WriteString("from offer_item_read_models oi ")
	sb.WriteString("left join offer_read_models o on o.offer_key = oi.offer_key ")
	sb.WriteString("left join offer_resource_read_models r on oi.resource_key = r.resource_key ")
	sb.WriteString("where ")
	sb.WriteString("   oi.from_user_key = ? ")
	sb.WriteString("or oi.to_user_key = ? ")
	sb.WriteString("or r.owner_user_key = ? ")

	if administeredGroupCount > 0 {
		sb.WriteString("or oi.from_group_key in " + groupsSb.String() + " ")
		sb.WriteString("or oi.to_group_key in " + groupsSb.String() + " ")
		sb.WriteString("or r.owner_group_key in " + groupsSb.String() + " ")
	}
	sb.WriteString("group by oi.offer_key")

	type result struct {
		OfferKey keys.OfferKey
	}
	var results []result

	err = q.db.Raw(sb.String(), params...).
		Find(&results).
		Error
	if err != nil {
		return nil, err
	}

	var offerKeys []keys.OfferKey
	for _, result := range results {
		offerKeys = append(offerKeys, result.OfferKey)
	}

	return keys.NewOfferKeys(offerKeys), nil

}
