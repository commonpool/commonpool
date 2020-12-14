package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/exceptions"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/mq"
)

func (g GroupService) CreateOrAcceptInvitation(ctx context.Context, request *group2.CreateOrAcceptInvitationRequest) (*group2.CreateOrAcceptInvitationResponse, error) {

	userSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	isNewMembership := false
	membershipKey := request.MembershipKey
	membershipToAccept, err := g.groupStore.GetMembership(ctx, membershipKey)
	if err != nil && !errors.Is(err, exceptions.ErrMembershipNotFound) {
		return nil, err
	}

	if err != nil && errors.Is(err, exceptions.ErrMembershipNotFound) {
		isNewMembership = true
	}

	if userSession.GetUserKey() == membershipKey.UserKey {

		if isNewMembership {
			newMembership, err := g.groupStore.CreateMembership(ctx, membershipKey, false, false, false, false, false, true)
			if err != nil {
				return nil, err
			}
			return &group2.CreateOrAcceptInvitationResponse{
				Membership: newMembership,
			}, nil
		}

		if membershipToAccept.UserConfirmed {
			err := fmt.Errorf("membership already confirmed")
			return nil, err
		}

		err := g.groupStore.MarkInvitationAsAccepted(ctx, membershipKey, group2.PartyUser)
		if err != nil {
			return nil, err
		}

	} else {

		loggedInUserMembershipKey := model.NewMembershipKey(membershipKey.GroupKey, userSession.GetUserKey())
		loggedInUserMembership, err := g.groupStore.GetMembership(ctx, loggedInUserMembershipKey)
		if err != nil {
			return nil, exceptions.ErrMembershipPartyUnauthorized
		}

		if !loggedInUserMembership.IsAdmin {
			return nil, exceptions.ErrManageMembershipsNotAdmin
		}

		if isNewMembership {
			newMembership, err := g.groupStore.CreateMembership(ctx, membershipKey, false, false, false, false, true, false)
			if err != nil {
				return nil, err
			}
			return &group2.CreateOrAcceptInvitationResponse{
				Membership: newMembership,
			}, nil
		}

		if membershipToAccept.GroupConfirmed {
			err := fmt.Errorf("already accepted")
			return nil, err
		}

		err = g.groupStore.MarkInvitationAsAccepted(ctx, membershipKey, group2.PartyGroup)
		if err != nil {
			return nil, err
		}

	}

	acceptedMembership, err := g.groupStore.GetMembership(ctx, membershipKey)
	if err != nil {
		return nil, err
	}

	if acceptedMembership.GroupConfirmed && acceptedMembership.UserConfirmed {
		// add user to group channel

		grp, err := g.groupStore.GetGroup(ctx, request.MembershipKey.GroupKey)
		if err != nil {
			return nil, err
		}

		usernameJoiningGroup, err := g.authStore.GetUsername(membershipKey.UserKey)
		if err != nil {
			return nil, err
		}

		channelSubscriptionKey := model.NewChannelSubscriptionKey(acceptedMembership.GetGroupKey().GetChannelKey(), acceptedMembership.GetUserKey())
		_, err = g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
		if err != nil {
			return nil, err
		}

		amqpChannel, err := g.amqpClient.GetChannel()
		if err != nil {
			return nil, err
		}

		channelKey := request.MembershipKey.GroupKey.GetChannelKey()
		err = amqpChannel.ExchangeBind(ctx, membershipKey.UserKey.GetExchangeName(), "", mq.WebsocketMessagesExchange, false, map[string]interface{}{
			"event_type": "chat.message",
			"channel_id": channelKey.String(),
			"x-match":    "all",
		})
		if err != nil {
			return nil, err
		}

		text := fmt.Sprintf("%s has joined #%s", usernameJoiningGroup, grp.Name)
		message := chat.NewContextBlock([]chat.BlockElement{
			chat.NewMarkdownObject(text)},
			nil,
		)

		_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(request.MembershipKey.GroupKey, membershipKey.UserKey, usernameJoiningGroup, text, []chat.Block{*message}, []chat.Attachment{}, nil))
		if err != nil {
			return nil, err
		}

	}

	return &group2.CreateOrAcceptInvitationResponse{
		Membership: acceptedMembership,
	}, nil

}
