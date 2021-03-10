package domain

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/exceptions"
	"time"
)

type ResourceValueType string

const (
	FromToDuration ResourceValueType = "from_to_duration"
)

func (r *ResourceValueType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	if str != string(FromToDuration) {
		return exceptions.ErrBadRequestf("invalid ResourceValueType '%s'", str)
	}
	*r = FromToDuration
	return nil
}

type ResourceValueEstimation struct {
	ValueType         ResourceValueType `json:"valueType" gorm:"not null;type:varchar(128)"`
	ValueFromDuration time.Duration     `json:"timeValueFrom" gorm:"not null"`
	ValueToDuration   time.Duration     `json:"timeValueTo" gorm:"not null"`
}

func (r ResourceValueEstimation) WithValueType(valueType ResourceValueType) ResourceValueEstimation {
	return ResourceValueEstimation{
		ValueType:         valueType,
		ValueFromDuration: r.ValueFromDuration,
		ValueToDuration:   r.ValueFromDuration,
	}
}

func (r ResourceValueEstimation) WithFromToValueType() ResourceValueEstimation {
	return r.WithValueType(FromToDuration)
}

func (r ResourceValueEstimation) WithValueFromDuration(from time.Duration) ResourceValueEstimation {
	return ResourceValueEstimation{
		ValueType:         r.ValueType,
		ValueFromDuration: from,
		ValueToDuration:   r.ValueFromDuration,
	}
}

func (r ResourceValueEstimation) WithHoursFromTo(fromHours, toHours int) ResourceValueEstimation {
	return r.
		WithValueFromDuration(time.Duration(fromHours) * time.Hour).
		WithValueToDuration(time.Duration(toHours) * time.Hour)
}

func (r ResourceValueEstimation) WithValueToDuration(to time.Duration) ResourceValueEstimation {
	return ResourceValueEstimation{
		ValueType:         r.ValueType,
		ValueFromDuration: r.ValueFromDuration,
		ValueToDuration:   to,
	}
}

func NewResourceValueEstimation() ResourceValueEstimation {
	return ResourceValueEstimation{}
}

type ResourceInfoBase struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	CallType     CallType     `json:"callType"`
	ResourceType ResourceType `json:"resourceType"`
}

type ResourceInfo struct {
	ResourceInfoBase
	Value ResourceValueEstimation `json:"value"`
}

func (r ResourceInfo) AsUpdate() ResourceInfoUpdate {
	return ResourceInfoUpdate{
		Name:        r.Name,
		Description: r.Description,
		Value:       r.Value,
	}
}

type ResourceInfoUpdate struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Value       ResourceValueEstimation `json:"value"`
}

func (r ResourceInfo) WithName(name string) ResourceInfo {
	return ResourceInfo{
		ResourceInfoBase: ResourceInfoBase{
			Name:         name,
			Description:  r.Description,
			CallType:     r.CallType,
			ResourceType: r.ResourceType,
		},
		Value: r.Value,
	}
}
func (r ResourceInfo) WithDescription(description string) ResourceInfo {
	return ResourceInfo{
		ResourceInfoBase: ResourceInfoBase{
			Name:         r.Name,
			Description:  description,
			CallType:     r.CallType,
			ResourceType: r.ResourceType,
		},
		Value: r.Value,
	}
}
func (r ResourceInfo) WithCallType(callType CallType) ResourceInfo {
	return ResourceInfo{
		ResourceInfoBase: ResourceInfoBase{
			Name:         r.Name,
			Description:  r.Description,
			CallType:     callType,
			ResourceType: r.ResourceType,
		},
		Value: r.Value,
	}
}

func (r ResourceInfo) WithIsOffer() ResourceInfo {
	return r.WithCallType(Offer)
}

func (r ResourceInfo) WithIsRequest() ResourceInfo {
	return r.WithCallType(Request)
}

func (r ResourceInfo) WithResourceType(resourceType ResourceType) ResourceInfo {
	return ResourceInfo{
		Value: r.Value,
		ResourceInfoBase: ResourceInfoBase{
			Name:         r.ResourceInfoBase.Name,
			Description:  r.ResourceInfoBase.Description,
			CallType:     r.ResourceInfoBase.CallType,
			ResourceType: resourceType,
		},
	}
}

func (r ResourceInfo) WithIsService() ResourceInfo {
	return r.WithResourceType(ServiceResource)
}

func (r ResourceInfo) WithIsObject() ResourceInfo {
	return r.WithResourceType(ObjectResource)
}

func (r ResourceInfo) WithValue(value ResourceValueEstimation) ResourceInfo {
	return ResourceInfo{
		ResourceInfoBase: ResourceInfoBase{
			Name:         r.Name,
			Description:  r.Description,
			CallType:     r.CallType,
			ResourceType: r.ResourceType,
		},
		Value: value,
	}
}

func NewResourceInfo() ResourceInfo {
	return ResourceInfo{}
}
