package group

import (
	"github.com/commonpool/backend/model"
	"github.com/satori/go.uuid"
)

func ParseGroupKey(value string) (model.GroupKey, error) {
	offerId, err := uuid.FromString(value)
	if err != nil {
		return model.GroupKey{}, ErrInvalidGroupId
	}
	return model.NewGroupKey(offerId), err
}
