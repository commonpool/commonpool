package groups

import (
	"cp/pkg/api"
	"cp/pkg/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Store interface {
	Create(group *api.Group) error
	Search() ([]*api.Group, error)
	Get(id string) (*api.Group, error)
	Update(group *api.Group) error
}

type GroupStore struct {
	db *gorm.DB
}

func (g *GroupStore) Search() ([]*api.Group, error) {
	var groups []*api.Group
	if err := g.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func NewGroupStore(db *gorm.DB) *GroupStore {
	return &GroupStore{
		db: db,
	}
}

func (g *GroupStore) Create(group *api.Group) error {
	return g.db.Create(group).Error
}

func (g *GroupStore) Get(id string) (*api.Group, error) {
	var group api.Group
	err := g.db.
		Preload("Memberships.User").
		Preload("Posts.Author").
		Model(&api.Group{}).First(&group, "ID = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, echo.ErrNotFound
	}
	if err := utils.CountMessages(g.db, group.Posts); err != nil {
		return nil, err
	}
	return &group, nil
}

func (g *GroupStore) Update(group *api.Group) error {
	err := g.db.Save(group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.ErrNotFound
	}
	return nil
}

var _ Store = &GroupStore{}
