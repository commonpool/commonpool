package domain

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/exceptions"
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

type ResourceInfoBase struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	CallType     CallType     `json:"callType"`
	ResourceType ResourceType `json:"resourceType"`
}

type ResourceInfo struct {
	ResourceInfoBase
}

func (r ResourceInfo) AsUpdate() ResourceInfoUpdate {
	return ResourceInfoUpdate{
		Name:        r.Name,
		Description: r.Description,
	}
}

type ResourceInfoUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r ResourceInfo) WithName(name string) ResourceInfo {
	return ResourceInfo{
		ResourceInfoBase: ResourceInfoBase{
			Name:         name,
			Description:  r.Description,
			CallType:     r.CallType,
			ResourceType: r.ResourceType,
		},
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

func NewResourceInfo() ResourceInfo {
	return ResourceInfo{}
}
