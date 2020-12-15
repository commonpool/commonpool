package model

import (
	"github.com/commonpool/backend/pkg/exceptions"
)

type Type int

const (
	Offer Type = iota
	Request
)

func ParseResourceType(s string) (*Type, error) {
	var res Type
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

	err := exceptions.ErrParseResourceType(s)
	return nil, &err
}
