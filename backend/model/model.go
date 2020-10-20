package model

import (
	"github.com/commonpool/backend/errors"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type ResourceType int

const (
	Offer ResourceType = iota
	Request
)

func ParseResourceType(s string) (*ResourceType, error) {
	var res ResourceType
	if s == "" {
		return nil, nil
	}
	if s == "0" {
		res = Offer
		return &res, nil
	}
	if s == "1" {
		res = Request
		return &res, nil
	}
	return nil, errors.NewError(errors.ErrInvalidResourceType, errors.ErrInvalidResourceTypeCode, http.StatusBadRequest)
}

type TimeSensitivity struct {
	Value int `gorm:"column:time_sensitivity;not null;"`
}
type NecessityLevel struct {
	Value int `gorm:"column:necessity_level;not null;"`
}
type ExchangeValue struct {
	Value int `gorm:"column:exchange_value;not null;"`
}

func NewTimeSensitivity(value int) TimeSensitivity {
	return TimeSensitivity{Value: value}
}

func NewNecessityLevel(value int) NecessityLevel {
	return NecessityLevel{Value: value}
}

func NewExchangeValue(value int) ExchangeValue {
	return ExchangeValue{Value: value}
}

type Resource struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time   `sql:"index"`
	Summary          string       `gorm:"not null;"`
	Description      string       `gorm:"not null;"`
	CreatedBy        string       `gorm:"not null;"`
	Type             ResourceType `gorm:"not null;"`
	ValueInHoursFrom int          `gorm:"not null;'"`
	ValueInHoursTo   int          `gorm:"not null"`
}

type User struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Username  string     `gorm:"not null"`
	Email     string     `gorm:"not null"`
}

type ResourceKey struct {
	uuid uuid.UUID
}

type UserKey struct {
	subject string
}

func NewUserKey(subject string) UserKey {
	return UserKey{subject: subject}
}

func (k *UserKey) String() string {
	return k.subject
}

func NewResourceKey() ResourceKey {
	return ResourceKey{
		uuid: uuid.NewV4(),
	}
}

func ParseResourceKey(key string) (*ResourceKey, error) {
	resourceUuid, err := uuid.FromString(key)
	if err != nil {
		return nil, err
	}
	resourceKey := ResourceKey{
		uuid: resourceUuid,
	}
	return &resourceKey, nil
}

func (r *ResourceKey) GetUUID() uuid.UUID {
	return r.uuid
}

func (r *ResourceKey) String() string {
	return r.uuid.String()
}

func (r *Resource) GetKey() ResourceKey {
	return ResourceKey{
		uuid: r.ID,
	}
}

func NewResource(
	key ResourceKey,
	resourceType ResourceType,
	createdBy string,
	summary string,
	description string,
	valueInHoursFrom int,
	valueInHoursTo int,
) Resource {
	return Resource{
		ID:               key.uuid,
		Summary:          summary,
		Description:      description,
		CreatedBy:        createdBy,
		Type:             resourceType,
		ValueInHoursFrom: valueInHoursFrom,
		ValueInHoursTo:   valueInHoursTo,
	}
}
