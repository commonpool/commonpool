package handler

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	domain2 "github.com/commonpool/backend/pkg/trading/domain"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type CreateResourceRequest struct {
	Resource CreateResourcePayload `json:"resource"`
}

type CreateResourcePayload struct {
	Summary          string                 `json:"summary" validate:"required,max=100"`
	Description      string                 `json:"description" validate:"required,max=2000"`
	Type             domain.ResourceType    `json:"resource_type"`
	CallType         domain.CallType        `json:"call_type"`
	ValueInHoursFrom int                    `json:"valueInHoursFrom" validate:"min=0"`
	ValueInHoursTo   int                    `json:"valueInHoursTo" validate:"min=0"`
	SharedWith       []InputResourceSharing `json:"sharedWith"`
}

type CreateResourceResponse struct {
	Resource Resource `json:"resource"`
}

// CreateResource
// @Summary Creates a resource
// @Description Creates a resource
// @ID createResource
// @Tags resources
// @Accept json
// @Produce json
// @Param resource body web.CreateResourceRequest true "Resource to create"
// @Success 200 {object} web.CreateResourceResponse
// @Failure 400 {object} utils.Error
// @Router /resources [post]
func (h *ResourceHandler) CreateResource(c echo.Context) error {

	var err error

	ctx, _ := handler.GetEchoContext(c, "CreateResource")

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	req := CreateResourceRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}
	if err = c.Validate(req); err != nil {
		return err
	}

	sharedWithGroupKeys, err, done := h.parseGroupKeys(c, req.Resource.SharedWith)
	if done {
		return err
	}

	resourceKey := keys.NewResourceKey(uuid.NewV4())
	resource := domain.NewResource(resourceKey)
	err = resource.Register(loggedInUser.GetUserKey(), *domain2.NewUserTarget(loggedInUser.GetUserKey()), domain.ResourceInfo{
		Value: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: time.Duration(req.Resource.ValueInHoursFrom) * time.Hour,
			ValueToDuration:   time.Duration(req.Resource.ValueInHoursTo) * time.Hour,
		},
		Name:         req.Resource.Summary,
		Description:  req.Resource.Description,
		CallType:     req.Resource.CallType,
		ResourceType: req.Resource.Type,
	}, *sharedWithGroupKeys)
	if err != nil {
		return err
	}

	var rm *readmodel.ResourceReadModel
	err = retry.Do(func() error {
		rm, err = h.getResource.Get(ctx, resourceKey)
		if err != nil {
			return err
		}
		if rm.Version != resource.GetVersion() {
			return fmt.Errorf("unexpected read model version")
		}
		return nil
	})
	if err != nil {
		return err
	}

	sharings, err := h.getResourceSharings.Get(ctx, resourceKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, CreateResourceResponse{
		Resource: NewResourceResponse(rm, sharings),
	})
}
