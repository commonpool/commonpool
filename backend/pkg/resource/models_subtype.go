package resource

import "github.com/commonpool/backend/errors"

type SubType string

const (
	ServiceResource SubType = "service"
	ObjectResource  SubType = "object"
)

func ParseResourceSubType(s string) (*SubType, error) {
	var res SubType
	if s == "" {
		return nil, nil
	}
	if s == "object" {
		res = ObjectResource
		return &res, nil
	}
	if s == "service" {
		res = ServiceResource
		return &res, nil
	}

	err := errors.ErrParseResourceType(s)
	return nil, &err
}
