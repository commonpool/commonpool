package resource

import "github.com/commonpool/backend/errors"

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

	err := errors.ErrParseResourceType(s)
	return nil, &err
}
