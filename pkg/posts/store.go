package posts

import (
	"cp/pkg/api"
	"cp/pkg/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type FindPostsOptions struct {
	IncludeDeleted bool
	Type           *api.PostType
	Query          *string
}

type Store interface {
	Create(post *api.Post) error
	Update(post *api.Post) error
	Get(postID string) (*api.Post, error)
	GetByAuthor(authorID string, options ...*FindPostsOptions) ([]*api.Post, error)
	GetByGroup(groupID string, options ...*FindPostsOptions) ([]*api.Post, error)
	Delete(postID string) error
}

type PostStore struct {
	db *gorm.DB
}

func (p *PostStore) Create(post *api.Post) error {
	return p.db.Create(post).Error
}

func (p *PostStore) Update(post *api.Post) error {
	return p.db.Save(post).Error
}

func (p *PostStore) Get(postID string) (*api.Post, error) {
	var result api.Post

	err := p.db.
		Preload("Author").
		Preload("Images").
		Preload("Group").
		First(&result, "id = ?", postID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, echo.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := utils.CountMessages(p.db, []*api.Post{&result}); err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *PostStore) GetByAuthor(authorID string, options ...*FindPostsOptions) ([]*api.Post, error) {
	var result []*api.Post
	db := p.db
	if len(options) > 0 && options[0].IncludeDeleted {
		db = db.Unscoped()
	}
	if err := db.
		Preload("Group").
		Preload("Author").
		Preload("Images").
		Model(&api.Post{}).
		Order("created_at desc").
		Find(&result, "author_id = ?", authorID).
		Error; err != nil {
		return nil, err
	}
	if err := utils.CountMessages(db, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostStore) GetByGroup(groupID string, options ...*FindPostsOptions) ([]*api.Post, error) {
	var result []*api.Post
	db := p.db
	if len(options) > 0 && options[0].IncludeDeleted {
		db = db.Unscoped()
	}
	query := db.
		Preload("Group").
		Preload("Author").
		Preload("Images").
		Model(&api.Post{})

	if len(options) > 0 {
		if options[0].Query != nil {
			qry := "%" + *options[0].Query + "%"
			query = query.Where("title like ? or description like ?", qry, qry)
		}
		if options[0].Type != nil {
			query = query.Where("type = ?", *options[0].Type)
		}
	}

	if err := query.
		Debug().
		Order("created_at desc").
		Find(&result, "group_id = ?", groupID).
		Error; err != nil {
		return nil, err
	}
	if err := utils.CountMessages(db, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostStore) Delete(postID string) error {
	return p.db.Delete(&api.Post{}, "id = ?", postID).Error
}

func NewPostStore(db *gorm.DB) *PostStore {
	return &PostStore{db: db}
}
