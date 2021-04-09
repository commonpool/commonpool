package images

import (
	"cp/pkg/api"
	"gorm.io/gorm"
)

type Store interface {
	Add(images []*api.Image) error
	Get(postID string) ([]*api.Image, error)
	Delete(imageID string) error
}

type ImageStore struct {
	db *gorm.DB
}

func (i *ImageStore) Add(images []*api.Image) error {
	if len(images) == 0 {
		return nil
	}
	return i.db.Create(images).Error
}

func (i *ImageStore) Delete(imageID string) error {
	return i.db.Delete(&api.Image{ID: imageID}).Error
}

func (i *ImageStore) Get(postID string) ([]*api.Image, error) {
	var result []*api.Image
	if err := i.db.Model(&api.Image{}).Where("post_id = ?", postID).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func NewImageStore(db *gorm.DB) *ImageStore {
	return &ImageStore{db: db}
}
