package handler

import (
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
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
	Type             resource.Type           `json:"type"`
	SubType          resource.SubType        `json:"subType"`
	CreatedAt        time.Time               `json:"createdAt"`
	CreatedBy        string                  `json:"createdBy"`
	CreatedById      string                  `json:"createdById"`
	ValueInHoursFrom int                     `json:"valueInHoursFrom"`
	ValueInHoursTo   int                     `json:"valueInHoursTo"`
	SharedWith       []OutputResourceSharing `json:"sharedWith"`
}

func NewResourceResponse(res *resource.Resource, creatorUsername string, creatorId string, sharedWithGroups *group2.Groups) Resource {

	//goland:noinspection GoPreferNilSlice
	var sharings = []OutputResourceSharing{}
	for _, withGroup := range sharedWithGroups.Items {
		sharings = append(sharings, OutputResourceSharing{
			GroupID:   withGroup.Key.String(),
			GroupName: withGroup.Name,
		})
	}

	return Resource{
		Id:               res.Key.String(),
		Type:             res.Type,
		SubType:          res.SubType,
		Description:      res.Description,
		Summary:          res.Summary,
		CreatedBy:        creatorUsername,
		CreatedById:      creatorId,
		CreatedAt:        res.CreatedAt,
		ValueInHoursFrom: res.ValueInHoursFrom,
		ValueInHoursTo:   res.ValueInHoursTo,
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
