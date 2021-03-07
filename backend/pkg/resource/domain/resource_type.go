package domain

import "github.com/commonpool/backend/pkg/exceptions"

type ResourceType string

const (
	ServiceResource ResourceType = "service"
	ObjectResource  ResourceType = "object"
)

func ParseResourceType(str string) (ResourceType, error) {
	switch str {
	case string(ServiceResource):
		return ServiceResource, nil
	case string(ObjectResource):
		return ObjectResource, nil
	default:
		return "", exceptions.ErrBadRequest("invalid resource type")
	}
}

type CallType string

const (
	Offer   CallType = "offer"
	Request CallType = "request"
)

func ParseCallType(str string) (CallType, error) {
	switch str {
	case string(Offer):
		return Offer, nil
	case string(Request):
		return Request, nil
	default:
		return "", exceptions.ErrBadRequest("invalid call type")
	}
}
