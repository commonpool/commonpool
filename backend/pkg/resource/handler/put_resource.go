package handler

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

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

// UpdateResource
// @Summary Updates a resource
// @Description Updates a resource
// @ID updateResource
// @Tags resources
// @Accept json
// @Produce json
// @Param id path string true "Resource id" format(uuid)
// @Param resource body web.UpdateResourceRequest true "Resource to create"
// @Success 200 {object} web.UpdateResourceResponse
// @Failure 400 {object} utils.Error
// @Router /resources [put]
func (h *ResourceHandler) UpdateResource(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "UpdateResource")
	l = l.Named("ResourceHandler.UpdateResource")

	l.Debug("getting logged in user")
	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	l.Debug("binding request")
	req := UpdateResourceRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	l.Debug("validating request")
	if err := c.Validate(req); err != nil {
		return err
	}

	l.Debug("parsing resource key")
	resourceKey, err := keys.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	l.Debug("loading resource from repo")
	resource, err := h.resourceRepo.Load(ctx, resourceKey)
	if err != nil {
		return err
	}
	if resource.GetVersion() == 0 {
		return exceptions.ErrResourceNotFound
	}

	l.Debug("changing resource info")
	err = resource.ChangeInfo(loggedInUser.GetUserKey(), domain.ResourceInfo{
		Value: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: time.Duration(req.Resource.ValueInHoursFrom) * time.Hour,
			ValueToDuration:   time.Duration(req.Resource.ValueInHoursTo) * time.Hour,
		},
		Name:         req.Resource.Summary,
		Description:  req.Resource.Description,
		CallType:     resource.GetCallType(),
		ResourceType: resource.GetResourceType(),
	})
	if err != nil {
		return err
	}

	l.Debug("getting resource readmodel")
	var rm *readmodel.ResourceReadModel
	err = retry.Do(func() error {
		rm, err = h.getResource.Get(ctx, resourceKey)
		if err != nil {
			return err
		}
		if rm.Version != resource.GetVersion() {
			l.Debug("read model version not up to date", zap.Int("expected", resource.GetVersion()), zap.Int("actual", rm.Version))
			return fmt.Errorf("unexpected read model version")
		}
		return nil
	})
	if err != nil {
		return err
	}

	l.Debug("getting sharings read model")
	sharings, err := h.getResourceSharings.Get(ctx, resourceKey)
	if err != nil {
		return err
	}

	l.Debug("returning response")
	return c.JSON(http.StatusOK, GetResourceResponse{
		Resource: NewResourceResponse(rm, sharings),
	})

}
