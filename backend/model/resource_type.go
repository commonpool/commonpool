package model

import (
	"github.com/commonpool/backend/errors"
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

	err := errors.ErrParseResourceType(s)
	return nil, &err
}
