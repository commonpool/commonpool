package echo

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/client"
	"github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	handler2 "github.com/commonpool/backend/pkg/resource/handler"
	tradinghandler "github.com/commonpool/backend/pkg/trading/handler"
	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
	"net/http"
)

var _ client.Client = &Client{}

func (e *Client) CreateGroup(ctx context.Context, createGroup *handler.CreateGroupRequest, output *handler.GetGroupResponse) error {
	return e.do(ctx, http.MethodPost, "/api/v1/groups", http.StatusCreated, createGroup, output)
}
func (e *Client) JoinGroup(ctx context.Context, membershipKey keys.MembershipKey) error {
	return e.do(ctx, http.MethodPost, "/api/v1/memberships", http.StatusAccepted, handler.CreateOrAcceptInvitationRequest{
		UserKey:  membershipKey.UserKey,
		GroupKey: membershipKey.GroupKey,
	}, nil)
}
func (e *Client) LeaveGroup(ctx context.Context, membershipKey keys.MembershipKey) error {
	return e.do(ctx, http.MethodDelete, "/api/v1/memberships", http.StatusAccepted, handler.CancelOrDeclineInvitationRequest{
		UserKey:  membershipKey.UserKey,
		GroupKey: membershipKey.GroupKey,
	}, nil)
}
func (e *Client) GetLoggedInUserMemberships(ctx context.Context, output *handler.GetMembershipsResponse) error {
	return e.do(ctx, http.MethodGet, "/api/v1/memberships", http.StatusOK, nil, output)
}
func (e *Client) GetMembership(ctx context.Context, membershipKey keys.MembershipKey, output *handler.GetMembershipResponse) error {
	return e.do(ctx, http.MethodGet, fmt.Sprintf("/api/v1/groups/%s/memberships/%s", membershipKey.GroupKey.String(), membershipKey.UserKey.String()), http.StatusOK, nil, output)
}
func (e *Client) GetMemberInvitationPicker(ctx context.Context, groupKey keys.GroupKey, query string, skip, take int, output *handler.GetUsersForGroupInvitePickerResponse) error {
	return e.do(ctx, http.MethodGet, fmt.Sprintf("/api/v1/groups/%s/invite-member-picker?query=%s&skip=%d&take=%d", groupKey.String(), query, skip, take), http.StatusOK, nil, output)
}
func (e *Client) SubmitOffer(ctx context.Context, offer *tradinghandler.SendOfferRequest, out *tradinghandler.GetOfferResponse) error {
	return e.do(ctx, http.MethodPost, "/api/v1/offers", http.StatusCreated, offer, out)
}
func (e *Client) GetOffer(ctx context.Context, offerKey keys.OfferKeyGetter, out *tradinghandler.GetOfferResponse) error {
	return e.do(ctx, http.MethodGet, fmt.Sprintf("/api/v1/offers/%s", offerKey.GetOfferKey().String()), http.StatusCreated, nil, out)
}
func (e *Client) AcceptOffer(ctx context.Context, offerKey keys.OfferKeyGetter) error {
	return e.do(ctx, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/accept", offerKey.GetOfferKey().String()), http.StatusAccepted, nil, nil)
}
func (e *Client) DeclineOffer(ctx context.Context, offerKey keys.OfferKeyGetter) error {
	return e.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/offers/%s/decline", offerKey.GetOfferKey().String()), http.StatusAccepted, nil, nil)
}
func (e *Client) ConfirmResourceGiven(ctx context.Context, offerKey keys.OfferKeyGetter, offerItemKey keys.OfferItemKey) error {
	return e.do(ctx, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/items/%s/confirm/resource-given", offerKey.GetOfferKey().String(), offerItemKey.String()), http.StatusAccepted, nil, nil)
}
func (e *Client) ConfirmServiceGiven(ctx context.Context, offerKey keys.OfferKeyGetter, offerItemKey keys.OfferItemKey) error {
	return e.do(ctx, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/items/%s/confirm/service-given", offerKey.GetOfferKey().String(), offerItemKey.String()), http.StatusAccepted, nil, nil)
}
func (e *Client) ConfirmResourceBorrowed(ctx context.Context, offerKey keys.OfferKeyGetter, offerItemKey keys.OfferItemKey) error {
	return e.do(ctx, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/items/%s/confirm/resource-borrowed", offerKey.GetOfferKey().String(), offerItemKey.String()), http.StatusAccepted, nil, nil)
}
func (e *Client) ConfirmResourceReturned(ctx context.Context, offerKey keys.OfferKeyGetter, offerItemKey keys.OfferItemKey) error {
	return e.do(ctx, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/items/%s/confirm/resource-returned", offerKey.GetOfferKey().String(), offerItemKey.String()), http.StatusAccepted, nil, nil)
}
func (e *Client) CreateResource(ctx context.Context, resource *handler2.CreateResourceRequest, out *handler2.GetResourceResponse) error {
	return e.do(ctx, http.MethodPost, "/api/v1/resources", http.StatusCreated, resource, out)
}
func (e *Client) UpdateResource(ctx context.Context, resourceKey keys.ResourceKeyGetter, resource *handler2.UpdateResourceRequest, out *handler2.GetResourceResponse) error {
	return e.do(ctx, http.MethodPut, fmt.Sprintf("/api/v1/resources/%s", resourceKey.GetResourceKey()), http.StatusOK, resource, out)
}
func (e *Client) GetResource(ctx context.Context, resourceKey keys.ResourceKeyGetter, out *handler2.GetResourceResponse) error {
	return e.do(ctx, http.MethodGet, fmt.Sprintf("/api/v1/resources/%s", resourceKey.GetResourceKey().String()), http.StatusOK, nil, out)
}
func (e *Client) SearchResources(ctx context.Context, query string, callType *domain.CallType, resourceType *domain.ResourceType, skip, take int, sharedWithGroup keys.GroupKeyGetter, output *handler2.SearchResourcesResponse) error {
	var sharedStr = ""
	if sharedWithGroup != nil {
		sharedStr = sharedWithGroup.GetGroupKey().String()
	}
	callTypeStr := ""
	if callType != nil {
		callTypeStr = string(*callType)
	}
	return e.do(ctx, http.MethodGet, fmt.Sprintf("/api/v1/resources?take=%d&skip=%d&query=%s&call=%s&group_id=%s", take, skip, query, callTypeStr, sharedStr), http.StatusOK, nil, output)
}

func (e *Client) GetWebsocketClient() (*websocket.Conn, error) {
	d := wstest.NewDialer(e.echo)
	header, err := e.authentication.GetRequestHeader()
	if err != nil {
		return nil, err
	}
	c, resp, err := d.Dial("ws://"+"whatever"+"/api/v1/ws", header)
	if err != nil {
		return nil, err
	}
	if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
		return nil, fmt.Errorf("resp.StatusCode = %q, want %q", got, want)
	}
	return c, nil
}
