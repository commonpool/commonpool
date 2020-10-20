package web

import (
	"github.com/commonpool/backend/model"
	"github.com/labstack/echo/v4"
	"time"
)

type SearchResourcesResponse struct {
	TotalCount int        `json:"totalCount"`
	Take       int        `json:"take"`
	Skip       int        `json:"skip"`
	Resources  []Resource `json:"resources"`
}

type GetResourceResponse struct {
	Resource Resource `json:"resource"`
}

type Resource struct {
	Id               string             `json:"id"`
	Summary          string             `json:"summary"`
	Description      string             `json:"description"`
	Type             model.ResourceType `json:"type"`
	CreatedAt        time.Time          `json:"createdAt"`
	CreatedBy        string             `json:"createdBy"`
	CreatedById      string             `json:"createdById"`
	ValueInHoursFrom int                `json:"valueInHoursFrom"`
	ValueInHoursTo   int                `json:"valueInHoursTo"`
}

type CreateResourceResponse struct {
	Resource Resource `json:"resource"`
}

type CreateResourceRequest struct {
	Resource CreateResourcePayload `json:"resource"`
}

type CreateResourcePayload struct {
	Summary          string             `json:"summary"`
	Description      string             `json:"description"`
	Type             model.ResourceType `json:"type"`
	ValueInHoursFrom int                `json:"valueInHoursFrom"`
	ValueInHoursTo   int                `json:"valueInHoursTo"`
}

type UpdateResourceRequest struct {
	Resource UpdateResourcePayload `json:"resource"`
}

type UpdateResourcePayload struct {
	Summary          string             `json:"summary"`
	Description      string             `json:"description"`
	Type             model.ResourceType `json:"type"`
	ValueInHoursFrom int                `json:"valueInHoursFrom"`
	ValueInHoursTo   int                `json:"valueInHoursTo"`
}

type UpdateResourceResponse struct {
	Resource Resource `json:"resource"`
}

type UserAuthResponse struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	Username        string `json:"username"`
	Id              string `json:"id"`
}

type UserInfoResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func (r *CreateResourceRequest) bind(c echo.Context) {

}
