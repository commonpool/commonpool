package model

import (
	"github.com/commonpool/backend/errors"
)

type ResourceType int

const (
	ResourceOffer ResourceType = iota
	ResourceRequest
)

func ParseResourceType(s string) (*ResourceType, error) {
	var res ResourceType
	if s == "" {
		return nil, nil
	}
	if s == "0" {
		res = ResourceOffer
		return &res, nil
	}
	if s == "1" {
		res = ResourceRequest
		return &res, nil
	}

	err := errors.ErrParseResourceType(s)
	return nil, &err
}
