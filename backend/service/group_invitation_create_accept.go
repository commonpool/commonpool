package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

func (g GroupService) CreateOrAcceptInvitation(ctx context.Context, request *group.CreateOrAcceptInvitationRequest) (*group.CreateOrAcceptInvitationResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "CreateOrAcceptInvitation")

	l = l.With(zap.Object("membership", request.MembershipKey))

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session")
		return nil, err
	}

	isNewMembership := false
	membershipKey := request.MembershipKey
	membershipToAccept, err := g.groupStore.GetMembership(ctx, membershipKey)
	if err != nil && !errors.Is(err, group.ErrMembershipNotFound) {
		l.Error("could not get membership to accept", zap.Error(err))
		return nil, err
	}

	if err != nil && errors.Is(err, group.ErrMembershipNotFound) {
		isNewMembership = true
	}

	if userSession.GetUserKey() == membershipKey.UserKey {

		if isNewMembership {
			newMembership, err := g.groupStore.CreateMembership(ctx, membershipKey, false, false, false, false, false, true)
			if err != nil {
				return nil, err
			}
			return &group.CreateOrAcceptInvitationResponse{
				Membership: newMembership,
			}, nil
		}

		l.Debug("user was invited by group. he has to accept the invitation")

		if membershipToAccept.UserConfirmed {
			err := fmt.Errorf("membership already confirmed")
			l.Error(err.Error())
			return nil, err
		}

		err := g.groupStore.MarkInvitationAsAccepted(ctx, membershipKey, group.UserParty)
		if err != nil {
			l.Error("could not mark invitation as accepted", zap.Error(err))
			return nil, err
		}

	} else {

		l.Debug("user asked the group to join. group admin has to confirm the invitation")

		loggedInUserMembershipKey := model.NewMembershipKey(membershipKey.GroupKey, userSession.GetUserKey())
		loggedInUserMembership, err := g.groupStore.GetMembership(ctx, loggedInUserMembershipKey)
		if err != nil {
			return nil, group.ErrMembershipPartyUnauthorized
		}

		if !loggedInUserMembership.IsAdmin {
			return nil, group.ErrManageMembershipsNotAdmin
		}

		if isNewMembership {
			newMembership, err := g.groupStore.CreateMembership(ctx, membershipKey, false, false, false, false, true, false)
			if err != nil {
				return nil, err
			}
			return &group.CreateOrAcceptInvitationResponse{
				Membership: newMembership,
			}, nil
		}

		if membershipToAccept.GroupConfirmed {
			err := fmt.Errorf("already accepted")
			l.Error(err.Error())
			return nil, err
		}

		err = g.groupStore.MarkInvitationAsAccepted(ctx, membershipKey, group.GroupParty)
		if err != nil {
			l.Error("could not mark invitation as accepted", zap.Error(err))
			return nil, err
		}

	}

	acceptedMembership, err := g.groupStore.GetMembership(ctx, membershipKey)
	if err != nil {
		l.Error("could not get accepted membership", zap.Error(err))
		return nil, err
	}

	if acceptedMembership.GroupConfirmed && acceptedMembership.UserConfirmed {
		// add user to group channel

		grp, err := g.groupStore.GetGroup(ctx, request.MembershipKey.GroupKey)
		if err != nil {
			l.Error("could not get group", zap.Error(err))
			return nil, err
		}

		usernameJoiningGroup, err := g.authStore.GetUsername(membershipKey.UserKey)
		if err != nil {
			l.Error("could not get username of user leaving group", zap.Error(err))
			return nil, err
		}

		channelSubscriptionKey := model.NewChannelSubscriptionKey(acceptedMembership.GetGroupKey().GetChannelKey(), acceptedMembership.GetUserKey())
		_, err = g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
		if err != nil {
			l.Error("could not subscribe to channel", zap.Error(err))
			return nil, err
		}

		amqpChannel, err := g.amqpClient.GetChannel()
		if err != nil {
			l.Error("could not gat amqp client", zap.Error(err))
			return nil, err
		}

		channelKey := request.MembershipKey.GroupKey.GetChannelKey()
		err = amqpChannel.ExchangeBind(ctx, membershipKey.UserKey.GetExchangeName(), "", amqp.WebsocketMessagesExchange, false, map[string]interface{}{
			"event_type": "chat.message",
			"channel_id": channelKey.String(),
			"x-match":    "all",
		})
		if err != nil {
			l.Error("could not unbind exchanges", zap.Error(err))
			return nil, err
		}

		text := fmt.Sprintf("%s has joined #%s", usernameJoiningGroup, grp.Name)
		message := chat.NewContextBlock([]chat.BlockElement{
			chat.NewMarkdownObject(text)},
			nil,
		)

		_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(request.MembershipKey.GroupKey, membershipKey.UserKey, usernameJoiningGroup, text, []chat.Block{*message}, []chat.Attachment{}, nil))
		if err != nil {
			l.Error("could not send user leaving message", zap.Error(err))
			return nil, err
		}

	}

	return &group.CreateOrAcceptInvitationResponse{
		Membership: acceptedMembership,
	}, nil

}
