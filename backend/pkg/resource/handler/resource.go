package handler

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type InputResourceSharing struct {
	GroupID string `json:"groupId" validate:"required,uuid"`
}

func (h *ResourceHandler) parseGroupKeys(c echo.Context, sharedWith []InputResourceSharing) (*keys.GroupKeys, error, bool) {
	sharedWithGroupKeys := make([]keys.GroupKey, len(sharedWith))
	for i := range sharedWith {
		groupKeyStr := sharedWith[i].GroupID
		groupKey, err := keys.ParseGroupKey(groupKeyStr)
		if err != nil {
			return nil, c.String(http.StatusBadRequest, "invalid group key : "+groupKeyStr), true
		}
		sharedWithGroupKeys[i] = groupKey
	}
	return keys.NewGroupKeys(sharedWithGroupKeys), nil, false
}

type Resource struct {
	Id               string                  `json:"id"`
	Summary          string                  `json:"summary"`
	Description      string                  `json:"description"`
	Type             string                  `json:"type"`
	SubType          string                  `json:"subType"`
	CreatedAt        time.Time               `json:"createdAt"`
	CreatedBy        string                  `json:"createdBy"`
	CreatedById      string                  `json:"createdById"`
	ValueInHoursFrom int                     `json:"valueInHoursFrom"`
	ValueInHoursTo   int                     `json:"valueInHoursTo"`
	SharedWith       []OutputResourceSharing `json:"sharedWith"`
}

func NewResourceResponse(res *readmodel.ResourceReadModel, shares []*readmodel.ResourceSharingReadModel) Resource {

	//goland:noinspection GoPreferNilSlice
	var sharings = []OutputResourceSharing{}
	for _, share := range shares {
		sharings = append(sharings, OutputResourceSharing{
			GroupID:   share.GroupKey,
			GroupName: share.GroupName,
		})
	}

	return Resource{
		Id:               res.ResourceKey,
		Type:             string(res.ResourceType),
		SubType:          string(res.CallType),
		Description:      res.Description,
		Summary:          res.ResourceName,
		CreatedBy:        res.CreatedByName,
		CreatedById:      res.CreatedBy,
		CreatedAt:        res.CreatedAt,
		ValueInHoursFrom: int(res.ValueFromDuration.Hours()),
		ValueInHoursTo:   int(res.ValueToDuration.Hours()),
		SharedWith:       sharings,
	}
}

type OutputResourceSharing struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

type GetResourceResponse struct {
	Resource Resource `json:"resource"`
}
