package client

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	resourcehandler "github.com/commonpool/backend/pkg/resource/handler"
	tradinghandler "github.com/commonpool/backend/pkg/trading/handler"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type Client interface {
	// Groups
	CreateGroup(ctx context.Context, createGroup *grouphandler.CreateGroupRequest, output *grouphandler.GetGroupResponse) error
	JoinGroup(ctx context.Context, membershipKey keys.MembershipKey) error
	LeaveGroup(ctx context.Context, membershipKey keys.MembershipKey) error
	GetLoggedInUserMemberships(ctx context.Context, output *grouphandler.GetMembershipsResponse) error
	GetMembership(ctx context.Context, membershipKey keys.MembershipKey, output *grouphandler.GetMembershipResponse) error
	GetMemberInvitationPicker(ctx context.Context, groupKey keys.GroupKey, query string, skip, take int, response *grouphandler.GetUsersForGroupInvitePickerResponse) error
	// Resources
	GetResource(ctx context.Context, resourceKey keys.ResourceKeyGetter, out *resourcehandler.GetResourceResponse) error
	CreateResource(ctx context.Context, resource *resourcehandler.CreateResourceRequest, out *resourcehandler.GetResourceResponse) error
	UpdateResource(ctx context.Context, resourceKey keys.ResourceKeyGetter, resource *resourcehandler.UpdateResourceRequest, out *resourcehandler.GetResourceResponse) error
	SearchResources(ctx context.Context, query string, callType *domain.CallType, resourceType *domain.ResourceType, skip, take int, sharedWithGroup keys.GroupKeyGetter, output *resourcehandler.SearchResourcesResponse) error
	// Trading
	SubmitOffer(ctx context.Context, offer *tradinghandler.SendOfferRequest, out *tradinghandler.GetOfferResponse) error
	GetOffer(ctx context.Context, offerKey keys.OfferKeyGetter, out *tradinghandler.GetOfferResponse) error
	AcceptOffer(ctx context.Context, offerKey keys.OfferKeyGetter) error
	DeclineOffer(ctx context.Context, offerKey keys.OfferKeyGetter) error
	ConfirmResourceGiven(ctx context.Context, offerKey keys.OfferKeyGetter, offerItemKey keys.OfferItemKey) error
	ConfirmServiceGiven(ctx context.Context, offerKey keys.OfferKeyGetter, offerItemKey keys.OfferItemKey) error
	ConfirmResourceBorrowed(ctx context.Context, offerKey keys.OfferKeyGetter, offerItemKey keys.OfferItemKey) error
	ConfirmResourceReturned(ctx context.Context, offerKey keys.OfferKeyGetter, offerItemKey keys.OfferItemKey) error
	// Websocket
	GetWebsocketClient() (*websocket.Conn, error)
}

type Authentication interface {
	Apply(request *http.Request) error
	GetRequestHeader() (http.Header, error)
}

type MockAuthentication struct {
	user *models.UserSession
}

func (m *MockAuthentication) GetRequestHeader() (http.Header, error) {
	result := http.Header{}
	err := m.applyHeader(result)
	return result, err
}

func (m *MockAuthentication) applyHeader(header http.Header) error {
	if m.user != nil {
		header.Set("X-Debug-Username", m.user.Username)
		header.Set("X-Debug-Email", m.user.Email)
		header.Set("X-Debug-User-Id", m.user.Subject)
		header.Set("X-Debug-Is-Authenticated", strconv.FormatBool(true))
	}
	return nil
}

func (m *MockAuthentication) Apply(request *http.Request) error {
	return m.applyHeader(request.Header)
}

var _ Authentication = &MockAuthentication{}

func NewMockAuthentication(userSession *models.UserSession) *MockAuthentication {
	return &MockAuthentication{
		user: userSession,
	}
}

func JsonContentType(c *http.Request) {
	c.Header.Set("Content-Type", "application/json")
}
