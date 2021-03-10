package domain

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/exceptions"
)

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
		return "", exceptions.ErrBadRequestf("invalid resource type: %s", str)
	}
}

type CallType string

func (c *CallType) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	return c.fromString(str)
}

func (c *CallType) UnmarshalParam(param string) error {
	return c.fromString(param)
}

func (c *CallType) fromString(str string) error {
	if str == string(Offer) {
		*c = Offer
	} else if str == string(Request) {
		*c = Request
	} else {
		return fmt.Errorf("invalid call type")
	}
	return nil
}

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
