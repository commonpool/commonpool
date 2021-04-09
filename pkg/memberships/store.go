package memberships

import (
	"cp/pkg/api"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strings"
)

type GetMembershipsOptions struct {
	HasPermission *api.MembershipPermission
	GroupID       *string
	UserID        *string
	Preload       []string
}

type Store interface {
	Get(groupID string, userID string) (*api.Membership, error)
	Create(membership *api.Membership) error
	Update(membership *api.Membership) error
	Delete(membership *api.Membership) error
	Find(out interface{}, option *GetMembershipsOptions) error
}

type MembershipStore struct {
	db *gorm.DB
}

func NewMembershipStore(db *gorm.DB) *MembershipStore {
	return &MembershipStore{
		db: db,
	}
}

var _ Store = &MembershipStore{}

func (m *MembershipStore) Get(groupID string, userID string) (*api.Membership, error) {
	var result api.Membership
	err := m.db.Preload("Group").Preload("User").Model(&api.Membership{}).First(&result, "group_id = ? and user_id = ?", groupID, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, echo.ErrNotFound
	}
	return &result, err
}

func (m *MembershipStore) Create(membership *api.Membership) error {
	return m.db.Create(membership).Error
}

func (m *MembershipStore) Update(membership *api.Membership) error {
	err := m.db.Save(membership).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *MembershipStore) Delete(membership *api.Membership) error {
	return m.db.Delete(&api.Membership{}, "user_id = ? and group_id = ?", membership.UserID, membership.GroupID).Error
}

func (m *MembershipStore) Find(out interface{}, option *GetMembershipsOptions) error {
	var clauses []string
	var params []interface{}

	if option.GroupID != nil {
		clauses = append(clauses, "group_id = ?")
		params = append(params, *option.GroupID)
	}

	if option.UserID != nil {
		clauses = append(clauses, "user_id = ?")
		params = append(params, *option.UserID)
	}

	if option.HasPermission != nil {
		hasPermission := *option.HasPermission
		if hasPermission == api.Owner {
			clauses = append(clauses, "permission = ?")
			params = append(params, api.Owner)
		} else if hasPermission == api.Admin {
			clauses = append(clauses, "permission in (?,?)")
			params = append(params, api.Owner, api.Admin)
		}
	}

	query := m.db
	for _, preload := range option.Preload {
		query = query.Preload(preload)
	}

	var sql = strings.Join(clauses, " AND ")
	return query.Debug().Model(&api.Membership{}).Where(sql, params...).Find(out).Error

}

func removeIndex(s []*api.Membership, index int) []*api.Membership {
	return append(s[:index], s[index+1:]...)
}
