package web

import (
	"github.com/commonpool/backend/pkg/resource/model"
	"time"
)

type InputResourceSharing struct {
	GroupID string `json:"groupId" validate:"required,uuid"`
}

type OutputResourceSharing struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

type Resource struct {
	Id               string                  `json:"id"`
	Summary          string                  `json:"summary"`
	Description      string                  `json:"description"`
	Type             model.Type              `json:"type"`
	SubType          model.SubType           `json:"subType"`
	CreatedAt        time.Time               `json:"createdAt"`
	CreatedBy        string                  `json:"createdBy"`
	CreatedById      string                  `json:"createdById"`
	ValueInHoursFrom int                     `json:"valueInHoursFrom"`
	ValueInHoursTo   int                     `json:"valueInHoursTo"`
	SharedWith       []OutputResourceSharing `json:"sharedWith"`
}

type SearchResourcesResponse struct {
	TotalCount int        `json:"totalCount"`
	Take       int        `json:"take"`
	Skip       int        `json:"skip"`
	Resources  []Resource `json:"resources"`
}

type GetResourceResponse struct {
	Resource Resource `json:"resource"`
}

type CreateResourceRequest struct {
	Resource CreateResourcePayload `json:"resource"`
}

type CreateResourcePayload struct {
	Summary          string                 `json:"summary" validate:"required,max=100"`
	Description      string                 `json:"description" validate:"required,max=2000"`
	Type             model.Type             `json:"type" validate:"min=0,max=1"`
	SubType          model.SubType          `json:"subType"`
	ValueInHoursFrom int                    `json:"valueInHoursFrom" validate:"min=0"`
	ValueInHoursTo   int                    `json:"valueInHoursTo" validate:"min=0"`
	SharedWith       []InputResourceSharing `json:"sharedWith"`
}

type CreateResourceResponse struct {
	Resource Resource `json:"resource"`
}

type UpdateResourceRequest struct {
	Resource UpdateResourcePayload `json:"resource"`
}

type UpdateResourcePayload struct {
	Summary          string                 `json:"summary" validate:"required,max=100"`
	Description      string                 `json:"description" validate:"required,max=2000"`
	ValueInHoursFrom int                    `json:"valueInHoursFrom" validate:"min=0"`
	ValueInHoursTo   int                    `json:"valueInHoursTo" validate:"min=0"`
	SharedWith       []InputResourceSharing `json:"sharedWith"`
}

type UpdateResourceResponse struct {
	Resource Resource `json:"resource"`
}
