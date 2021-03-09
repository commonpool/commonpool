package test

import (
	resourcedomain "github.com/commonpool/backend/pkg/resource/domain"
	"time"
)

func AResourceInfo() resourcedomain.ResourceInfo {
	return resourcedomain.ResourceInfo{
		ResourceInfoBase: resourcedomain.ResourceInfoBase{
			Name:         "name",
			Description:  "description",
			CallType:     resourcedomain.Offer,
			ResourceType: resourcedomain.ServiceResource,
		},
		Value: resourcedomain.ResourceValueEstimation{
			ValueType:         resourcedomain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   4 * time.Hour,
		},
	}
}
