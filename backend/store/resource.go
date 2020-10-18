package store

import (
	"errors"
	"fmt"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type ResourceStore struct {
	db *gorm.DB
}

// GetByKey Gets a resource by keys
func (rs *ResourceStore) GetByKey(key model.ResourceKey, r *model.Resource) error {
	if err := rs.db.First(r, "id = ?", key.GetUUID().String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewResourceNotFoundError(key.String())
		}
		return err
	}
	return nil
}

// Search search for resources
func (rs *ResourceStore) Search(query resource.Query) (*resource.QueryResult, error) {
	var resources []model.Resource

	chain := rs.db.Model(resources)

	if query.Type != nil {
		chain = chain.Where(`"resources"."type" = ?`, query.Type)
	}

	if query.Query != nil {
		chain = chain.Where(`"resources"."summary" like ?`, "%"+*query.Query+"%")
	}

	var totalCount int64
	err := chain.Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	err = chain.
		Offset(query.Skip).
		Limit(query.Take).
		Order("created_at desc").
		Find(&resources).
		Error

	if err != nil {
		return nil, err
	}

	return &resource.QueryResult{
		Items:      resources,
		TotalCount: int(totalCount),
		Take:       query.Take,
		Skip:       query.Skip,
	}, nil

}

// Delete deletes a resource
func (rs *ResourceStore) Delete(key model.ResourceKey) error {
	result := rs.db.Delete(&model.Resource{}, "id = ?", key.GetUUID().String())
	if result.RowsAffected == 0 {
		return errs.NewResourceNotFoundError(key.String())
	}
	return result.Error
}

// Create creates a resource
func (rs *ResourceStore) Create(resource *model.Resource) error {
	return rs.db.Create(&resource).Error
}

// Update updates a resource
func (rs *ResourceStore) Update(resource *model.Resource) error {
	update := rs.db.Model(resource).Save(resource)
	if update.RowsAffected == 0 {
		key := resource.GetKey()
		return errs.NewResourceNotFoundError(key.String())
	}
	return update.Error
}

var _ resource.Store = &ResourceStore{}

func NewResourceStore(db *gorm.DB) *ResourceStore {
	return &ResourceStore{
		db: db,
	}
}

func NewTestDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./realworld_test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(3)
	return db
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Resource{},
	)
}

func DropTestDB() error {
	if err := os.Remove("./realworld_test.db"); err != nil {
		return err
	}
	return nil
}
