package domain

import "time"

type ResourceValueType string

const (
	FromToDuration ResourceValueType = "from_to_duration"
)

type ResourceValueEstimation struct {
	ValueType         ResourceValueType `json:"value_type" gorm:"not null;type:varchar(128)"`
	ValueFromDuration time.Duration     `json:"time_value_from" gorm:"not null"`
	ValueToDuration   time.Duration     `json:"time_value_to" gorm:"not null"`
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

type ResourceInfo struct {
	Value        ResourceValueEstimation `json:"value"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	CallType     CallType                `json:"call_type"`
	ResourceType ResourceType            `json:"resource_type"`
}

func (r ResourceInfo) WithName(name string) ResourceInfo {
	return ResourceInfo{
		Value:        r.Value,
		Name:         name,
		Description:  r.Description,
		CallType:     r.CallType,
		ResourceType: r.ResourceType,
	}
}
func (r ResourceInfo) WithDescription(description string) ResourceInfo {
	return ResourceInfo{
		Value:        r.Value,
		Name:         r.Name,
		Description:  description,
		CallType:     r.CallType,
		ResourceType: r.ResourceType,
	}
}
func (r ResourceInfo) WithCallType(callType CallType) ResourceInfo {
	return ResourceInfo{
		Value:        r.Value,
		Name:         r.Name,
		Description:  r.Description,
		CallType:     callType,
		ResourceType: r.ResourceType,
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
		Value:        r.Value,
		Name:         r.Name,
		Description:  r.Description,
		CallType:     r.CallType,
		ResourceType: resourceType,
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
		Value:        value,
		Name:         r.Name,
		Description:  r.Description,
		CallType:     r.CallType,
		ResourceType: r.ResourceType,
	}
}

func NewResourceInfo() ResourceInfo {
	return ResourceInfo{}
}
