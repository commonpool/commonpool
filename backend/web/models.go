package web

import (
	"github.com/labstack/echo/v4"
	"github.com/commonpool/backend/model"
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
	Id              string             `json:"id"`
	Summary         string             `json:"summary"`
	Description     string             `json:"description"`
	Type            model.ResourceType `json:"type"`
	TimeSensitivity int                `json:"timeSensitivity"`
	ExchangeValue   int                `json:"exchangeValue"`
	NecessityLevel  int                `json:"necessityLevel"`
}

type CreateResourceResponse struct {
	Resource Resource `json:"resource"`
}

type CreateResourceRequest struct {
	Resource CreateResourcePayload `json:"resource"`
}

type CreateResourcePayload struct {
	Summary         string             `json:"summary"`
	Description     string             `json:"description"`
	Type            model.ResourceType `json:"type"`
	TimeSensitivity int                `json:"timeSensitivity"`
	ExchangeValue   int                `json:"exchangeValue"`
	NecessityLevel  int                `json:"necessityLevel"`
}

type UpdateResourceRequest struct {
	Resource UpdateResourcePayload `json:"resource"`
}

type UpdateResourcePayload struct {
	Summary         string             `json:"summary"`
	Description     string             `json:"description"`
	Type            model.ResourceType `json:"type"`
	TimeSensitivity int                `json:"timeSensitivity"`
	ExchangeValue   int                `json:"exchangeValue"`
	NecessityLevel  int                `json:"necessityLevel"`
}

type UpdateResourceResponse struct {
	Resource Resource `json:"resource"`
}

func (r *CreateResourceRequest) bind(c echo.Context) {

}
