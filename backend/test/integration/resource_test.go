package integration

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/web"
	"github.com/go-playground/assert/v2"
	"net/http"
	"testing"
)

func CreateResource(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.CreateResourceRequest) (*web.CreateResourceResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/resources", request)
	PanicIfError(a.CreateResource(c))
	response := &web.CreateResourceResponse{}
	return response, ReadResponse(t, recorder, response)
}

func TestUserCanCreateResource(t *testing.T) {
	teardown()
	setup()

	ctx := context.Background()

	resp, httpResp := CreateResource(t, ctx, User1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.ResourceOffer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	assert.Equal(t, http.StatusCreated, httpResp.StatusCode)

	assert.Equal(t, "Summary", resp.Resource.Summary)
	assert.Equal(t, "Description", resp.Resource.Description)
	assert.Equal(t, resource.ResourceOffer, resp.Resource.Type)
	assert.Equal(t, 1, resp.Resource.ValueInHoursFrom)
	assert.Equal(t, 3, resp.Resource.ValueInHoursTo)

}
