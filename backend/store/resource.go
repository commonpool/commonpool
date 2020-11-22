package store

import (
	"errors"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/utils"
	"gorm.io/gorm"
	"strings"
)

type ResourceStore struct {
	db *gorm.DB
}

var _ resource.Store = &ResourceStore{}

func NewResourceStore(db *gorm.DB) *ResourceStore {
	return &ResourceStore{
		db: db,
	}
}

// GetByKey Gets a resource by key
func (rs *ResourceStore) GetByKey(getResourceByKeyQuery *resource.GetResourceByKeyQuery) *resource.GetResourceByKeyResponse {

	result := model.Resource{}
	resourceKey := getResourceByKeyQuery.ResourceKey
	resourceKeyStr := resourceKey.String()

	if err := rs.db.First(&result, "id = ?", resourceKeyStr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := errs.ErrResourceNotFound(resourceKeyStr)
			return resource.NewGetResourceByKeyResponseError(&response)
		}
		return resource.NewGetResourceByKeyResponseError(err)
	}

	sharings, err := getResourceSharings(rs.db, []model.ResourceKey{resourceKey})
	if err != nil {
		return resource.NewGetResourceByKeyResponseError(err)
	}

	return resource.NewGetResourceByKeyResponseSuccess(&result, sharings)
}

// Delete deletes a resource
func (rs *ResourceStore) Delete(deleteResourceQuery *resource.DeleteResourceQuery) *resource.DeleteResourceResponse {
	err := rs.db.Transaction(func(tx *gorm.DB) error {

		resourceKey := deleteResourceQuery.ResourceKey
		resourceKeyStr := resourceKey.GetUUID().String()

		err := rs.db.Delete(&model.ResourceSharing{}, "resource_id = ?", resourceKeyStr).Error
		if err != nil {
			return err
		}

		result := rs.db.Delete(&model.Resource{}, "id = ?", resourceKeyStr)
		if result.RowsAffected == 0 {
			response := errs.ErrResourceNotFound(resourceKeyStr)
			return &response
		}

		if result.Error != nil {
			return result.Error
		}

		return rs.db.Delete(&model.ResourceSharing{}, "resource_id = ?", resourceKey.String()).Error
	})
	return resource.NewDeleteResourceResponse(err)
}

// Create creates a resource
func (rs *ResourceStore) Create(createResourceQuery *resource.CreateResourceQuery) *resource.CreateResourceResponse {
	err := rs.db.Transaction(func(tx *gorm.DB) error {
		res := createResourceQuery.Resource

		err := tx.Create(res).Error
		if err != nil {
			return err
		}
		return createResourceSharings(createResourceQuery.SharedWith, res.GetKey(), tx)
	})
	return resource.NewCreateResourceResponse(err)
}

// Update updates a resource
func (rs *ResourceStore) Update(updateResourceRequest *resource.UpdateResourceQuery) *resource.UpdateResourceResponse {
	err := rs.db.Transaction(func(tx *gorm.DB) error {
		res := updateResourceRequest.Resource
		resKey := res.GetKey()

		update := tx.Model(res).Save(res)

		if update.RowsAffected == 0 {
			key := resKey
			response := errs.ErrResourceNotFound(key.String())
			return &response
		}

		if update.Error != nil {
			return update.Error
		}

		err := tx.Delete(&model.ResourceSharing{}, "resource_id = ?", res.ID.String()).Error
		if err != nil {
			return err
		}

		return createResourceSharings(updateResourceRequest.SharedWith, resKey, tx)
	})

	return resource.NewUpdateResourceResponse(err)

}

// Search search for resources
func (rs *ResourceStore) Search(query *resource.SearchResourcesQuery) *resource.SearchResourcesResponse {
	var resources []model.Resource

	chain := rs.db.Model(&model.Resource{})

	if query.Type != nil {
		chain = chain.Where(`"resources"."type" = ?`, query.Type)
	}

	if query.Query != nil {
		chain = chain.Where(`"resources"."summary" like ?`, "%"+*query.Query+"%")
	}

	if query.CreatedBy != "" {
		chain = chain.Where(`"resources"."created_by" = ?`, query.CreatedBy)
	}

	if query.SharedWithGroup != nil {
		println(query.SharedWithGroup.ID.String())
		chain = chain.Joins("join resource_sharings on (resource_sharings.resource_id = resources.id and resource_sharings.group_id = ?)", query.SharedWithGroup.ID.String())
	}

	var totalCount int64
	err := chain.Count(&totalCount).Error
	if err != nil {
		return resource.NewSearchResourcesResponseError(err)
	}

	err = chain.
		Offset(query.Skip).
		Limit(query.Take).
		Order("created_at desc").
		Find(&resources).
		Error

	if err != nil {
		return resource.NewSearchResourcesResponseError(err)
	}

	resourceKeys := make([]model.ResourceKey, len(resources))
	for i := range resources {
		resourceKeys[i] = resources[i].GetKey()
	}

	sharings, err := getResourceSharings(rs.db, resourceKeys)
	if err != nil {
		return resource.NewSearchResourcesResponseError(err)
	}

	return resource.NewSearchResourcesResponseSuccess(
		model.NewResources(resources),
		sharings,
		int(totalCount),
		query.Take,
		query.Skip)

}

func getResourceSharings(db *gorm.DB, resources []model.ResourceKey) (*model.ResourceSharings, error) {

	var sharings []model.ResourceSharing

	err := utils.Partition(len(resources), 999, func(i1 int, i2 int) error {
		var qryPart []string
		var qryParam []interface{}
		for _, res := range resources[i1:i2] {
			qryPart = append(qryPart, "?")
			qryParam = append(qryParam, res.String())
		}
		qry := "resource_id IN ( " + strings.Join(qryPart, ",") + ")"
		var part []model.ResourceSharing
		err := db.Model(model.ResourceSharing{}).Where(qry, qryParam...).Find(&part).Error
		if err != nil {
			return err
		}
		for _, sharing := range part {
			sharings = append(sharings, sharing)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return model.NewResourceSharings(sharings)
}

func createResourceSharings(with []model.GroupKey, resourceKey model.ResourceKey, db *gorm.DB) error {
	if len(with) == 0 {
		return nil
	}
	var resourceSharings = make([]model.ResourceSharing, len(with))
	for i, groupKey := range with {
		resourceSharing := model.NewResourceSharing(resourceKey, groupKey)
		resourceSharings[i] = resourceSharing
	}

	return db.Create(resourceSharings).Error
}
