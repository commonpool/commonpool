package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

var _ group.Service = &GroupService{}

type GroupService struct {
	groupStore  group.Store
	amqpClient  amqp.AmqpClient
	chatService chat.Service
	authStore   auth.Store
}

func NewGroupService(groupStore group.Store, amqpClient amqp.AmqpClient, chatService chat.Service, authStore auth.Store) *GroupService {
	return &GroupService{
		groupStore:  groupStore,
		amqpClient:  amqpClient,
		chatService: chatService,
		authStore:   authStore,
	}
}

func (g GroupService) GetGroupsByKeys(ctx context.Context, groupKeys []model.GroupKey) (*group.Groups, error) {

	ctx, l := GetCtx(ctx, "GroupService", "GetGroupsByKeys")

	groups, err := g.groupStore.GetGroupsByKeys(ctx, groupKeys)
	if err != nil {
		l.Error("could not get groups", zap.Error(err))
		return nil, err
	}

	return groups, nil
}

func (g GroupService) CreateGroup(ctx context.Context, request *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "CreateGroup")

	l.Debug("getting user session")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session", zap.Error(err))
		return nil, err
	}

	l.Debug("creating group")

	grp, err := g.groupStore.CreateGroup(ctx, request.GroupKey, userSession.GetUserKey(), request.Name, request.Description)
	if err != nil {
		l.Error("could not create group", zap.Error(err))
		return nil, err
	}

	l.Debug("creating membership for owner")

	membershipKey := model.NewMembershipKey(grp.GetKey(), userSession.GetUserKey())
	membership, err := g.groupStore.CreateMembership(ctx, membershipKey, true, true, true, false, true, true)
	if err != nil {
		l.Error("could not create membership", zap.Error(err))
		return nil, err
	}

	l.Debug("creating group channel")

	channel, err := g.chatService.CreateChannel(ctx, request.GroupKey.GetChannelKey(), chat.GroupChannel)
	if err != nil {
		l.Error("could not create channel for group", zap.Error(err))
		return nil, err
	}

	l.Debug("subscribing owner to group channel")

	channelSubscriptionKey := model.NewChannelSubscriptionKey(channel.GetKey(), userSession.GetUserKey())
	channelSubscription, err := g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
	if err != nil {
		l.Error("could not subscribe owner to group channel", zap.Error(err))
		return nil, err
	}

	return &group.CreateGroupResponse{
		ChannelKey:      channel.GetKey(),
		SubscriptionKey: channelSubscription.GetKey(),
		Group:           grp,
		Membership:      membership,
	}, nil

}

func (g GroupService) GetGroup(ctx context.Context, request *group.GetGroupRequest) (*group.GetGroupResult, error) {

	ctx, l := GetCtx(ctx, "GroupService", "GetGroup")
	l = l.With(zap.Object("group", request.Key))

	l.Debug("getting group")

	grp, err := g.groupStore.GetGroup(ctx, request.Key)
	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return nil, err
	}

	return &group.GetGroupResult{
		Group: grp,
	}, nil

}

func (g GroupService) GetGroups(ctx context.Context, request *group.GetGroupsRequest) (*group.GetGroupsResult, error) {

	ctx, l := GetCtx(ctx, "GroupService", "GetGroups")

	l.Debug("getting groups")

	groups, totalCount, err := g.groupStore.GetGroups(request.Take, request.Skip)
	if err != nil {
		l.Error("could not get groups", zap.Error(err))
		return nil, err
	}

	return &group.GetGroupsResult{
		Items:      groups,
		TotalCount: totalCount,
	}, nil

}

func (g GroupService) GetMembership(ctx context.Context, request *group.GetMembershipRequest) (*group.GetMembershipResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "GetMembership")
	l = l.With(zap.Object("membership", request.MembershipKey))

	l.Debug("getting membership")

	membership, err := g.groupStore.GetMembership(ctx, request.MembershipKey)
	if err != nil {
		l.Error("could not get membership", zap.Error(err))
		return nil, err
	}

	return &group.GetMembershipResponse{
		Membership: membership,
	}, nil

}

func (g GroupService) SendGroupInvitation(ctx context.Context, request *group.InviteRequest) (*group.InviteResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "SendGroupInvitation")
	membershipKey := request.MembershipKey

	l = l.With(zap.Object("memnership", membershipKey))

	l.Debug("sending group invitation")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session")
		return nil, err
	}

	var membership *group.Membership

	if userSession.GetUserKey() == request.MembershipKey.UserKey {

		// user invites himself to the group

		membership, err = g.groupStore.CreateMembership(ctx, request.MembershipKey, false, false, false, false, false, true)
		if err != nil {
			l.Error("could not create membership", zap.Error(err))
			return nil, err
		}

	} else {

		adminMembershipKey := model.NewMembershipKey(request.MembershipKey.GroupKey, userSession.GetUserKey())

		adminMembership, err := g.groupStore.GetMembership(ctx, adminMembershipKey)
		if err != nil {
			l.Error("could not get admin membership", zap.Error(err))
			return nil, err
		}

		if !adminMembership.IsAdmin {
			err := fmt.Errorf("cannot invite user if not admin")
			l.Error(err.Error())
			return nil, err
		}

		membership, err = g.groupStore.CreateMembership(ctx, request.MembershipKey, false, false, false, false, true, false)
		if err != nil {
			l.Error("could not create membership", zap.Error(err))
			return nil, err
		}

	}

	return &group.InviteResponse{
		Membership: membership,
	}, nil

}

func (g GroupService) GetUserMemberships(ctx context.Context, request *group.GetMembershipsForUserRequest) (*group.GetMembershipsForUserResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "GetUserMemberships")

	l = l.With(zap.Object("user", request.UserKey))

	l.Debug("getting user memberships")

	memberships, err := g.groupStore.GetMembershipsForUser(ctx, request.UserKey, request.MembershipStatus)
	if err != nil {
		l.Error("could not get memberships for user", zap.Error(err))
		return nil, err
	}

	return &group.GetMembershipsForUserResponse{
		Memberships: memberships,
	}, nil
}

func (g GroupService) GetGroupsMemberships(ctx context.Context, request *group.GetMembershipsForGroupRequest) (*group.GetMembershipsForGroupResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "GetGroupsMemberships")

	l = l.With(zap.Object("group", request.GroupKey))

	l.Debug("getting group memberships")

	memberships, err := g.groupStore.GetMembershipsForGroup(ctx, request.GroupKey, request.MembershipStatus)
	if err != nil {
		l.Error("could not get memberships for group", zap.Error(err))
		return nil, err
	}

	return &group.GetMembershipsForGroupResponse{
		Memberships: memberships,
	}, nil
}

func (g GroupService) AcceptInvitation(ctx context.Context, request *group.AcceptInvitationRequest) (*group.AcceptInvitationResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "AcceptInvitation")

	l = l.With(zap.Object("membership", request.MembershipKey))

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session")
		return nil, err
	}

	membershipToAcceptKey := request.MembershipKey
	membershipToAccept, err := g.groupStore.GetMembership(ctx, membershipToAcceptKey)
	if err != nil {
		l.Error("could not get membership to accept", zap.Error(err))
		return nil, err
	}

	if userSession.GetUserKey() == membershipToAcceptKey.UserKey {

		l.Debug("user was invited by group. he has to accept the invitation")

		if membershipToAccept.UserConfirmed {
			err := fmt.Errorf("membership already confirmed")
			l.Error(err.Error())
			return nil, err
		}

		err := g.groupStore.MarkInvitationAsAccepted(ctx, membershipToAcceptKey, group.UserParty)
		if err != nil {
			l.Error("could not mark invitation as accepted", zap.Error(err))
			return nil, err
		}

	} else {

		l.Debug("user asked the group to join. group admin has to confirm the invitation")

		loggedInUserMembershipKey := model.NewMembershipKey(membershipToAcceptKey.GroupKey, userSession.GetUserKey())
		loggedInUserMembership, err := g.groupStore.GetMembership(ctx, loggedInUserMembershipKey)
		if err != nil {
			l.Error("could not get membership for logged in user", zap.Error(err))
			return nil, err
		}

		if !loggedInUserMembership.IsAdmin {
			err := fmt.Errorf("cannot accept membership, not admin")
			l.Error(err.Error())
			return nil, err
		}

		if membershipToAccept.GroupConfirmed {
			err := fmt.Errorf("already accepted")
			l.Error(err.Error())
			return nil, err
		}

		err = g.groupStore.MarkInvitationAsAccepted(ctx, membershipToAcceptKey, group.GroupParty)
		if err != nil {
			l.Error("could not mark invitation as accepted", zap.Error(err))
			return nil, err
		}

	}

	acceptedMembership, err := g.groupStore.GetMembership(ctx, membershipToAcceptKey)
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

		usernameJoiningGroup, err := g.authStore.GetUsername(membershipToAcceptKey.UserKey)
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
		err = amqpChannel.ExchangeBind(ctx, membershipToAcceptKey.UserKey.GetExchangeName(), "", amqp.WebsocketMessagesExchange, false, map[string]interface{}{
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

		_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(request.MembershipKey.GroupKey, membershipToAcceptKey.UserKey, usernameJoiningGroup, text, []chat.Block{*message}, []chat.Attachment{}, nil))
		if err != nil {
			l.Error("could not send user leaving message", zap.Error(err))
			return nil, err
		}

	}

	return &group.AcceptInvitationResponse{
		Membership: acceptedMembership,
	}, nil

}

func (g GroupService) DeclineInvitation(ctx context.Context, request *group.DeclineInvitationRequest) error {

	ctx, l := GetCtx(ctx, "GroupService", "DeclineInvitation")

	membershipKeyToDecline := request.MembershipKey
	l = l.With(zap.Object("membership", membershipKeyToDecline))

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session")
		return err
	}

	membership, err := g.groupStore.GetMembership(ctx, membershipKeyToDecline)
	if err != nil {
		l.Error("could not get membership", zap.Error(err))
		return err
	}

	if membership.GetUserKey() == userSession.GetUserKey() {
		// user is declining invitaiton from group

		// todo: check if last membership

		err = g.groupStore.DeleteMembership(ctx, membershipKeyToDecline)
		if err != nil {
			l.Error("could not delete membership", zap.Error(err))
			return err
		}

	} else {
		// group is declining invitation from user

		adminMembershipKey := model.NewMembershipKey(membershipKeyToDecline.GroupKey, userSession.GetUserKey())
		adminMembership, err := g.groupStore.GetMembership(ctx, adminMembershipKey)
		if err != nil {
			l.Error("could not get admin membership", zap.Error(err))
			return err
		}

		if !adminMembership.IsAdmin {
			err := fmt.Errorf("cannot decline invitation if not admin")
			l.Error(err.Error())
			return err
		}

		err = g.groupStore.DeleteMembership(ctx, membershipKeyToDecline)
		if err != nil {
			l.Error("could not delete membership", zap.Error(err))
			return err
		}
	}

	return nil

}

func (g GroupService) LeaveGroup(ctx context.Context, request *group.LeaveGroupRequest) error {

	ctx, l := GetCtx(ctx, "GroupService", "LeaveGroup")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session")
		return err
	}

	membershipToDelete := request.MembershipKey
	l = l.With(zap.Object("membership", membershipToDelete))

	membership, err := g.groupStore.GetMembership(ctx, membershipToDelete)
	if err != nil {
		l.Error("could not get membership", zap.Error(err))
		return err
	}

	usernameLeavingGroup, err := g.authStore.GetUsername(membershipToDelete.UserKey)
	if err != nil {
		l.Error("could not get username of user leaving group", zap.Error(err))
		return err
	}

	if membership.GetUserKey() == userSession.GetUserKey() {

		// user is leaving group
		// todo: check if last membership

		err = g.groupStore.DeleteMembership(ctx, membershipToDelete)
		if err != nil {
			l.Error("could not delete membership", zap.Error(err))
			return err
		}

	} else {

		// group is kicking user

		adminMembershipKey := model.NewMembershipKey(membershipToDelete.GroupKey, userSession.GetUserKey())
		adminMembership, err := g.groupStore.GetMembership(ctx, adminMembershipKey)
		if err != nil {
			l.Error("could not get admin membership", zap.Error(err))
			return err
		}

		if !adminMembership.IsAdmin {
			err := fmt.Errorf("cannot delete membership if not admin")
			l.Error(err.Error())
			return err
		}

		err = g.groupStore.DeleteMembership(ctx, membershipToDelete)
		if err != nil {
			l.Error("could not delete membership", zap.Error(err))
			return err
		}
	}

	amqpChannel, err := g.amqpClient.GetChannel()
	if err != nil {
		l.Error("could not gat amqp client", zap.Error(err))
		return err
	}

	channelSubscriptionKey := model.NewChannelSubscriptionKey(membershipToDelete.GroupKey.GetChannelKey(), membershipToDelete.UserKey)
	err = g.chatService.UnsubscribeFromChannel(ctx, channelSubscriptionKey)
	if err != nil {
		l.Error("could not unsubscribe to channel", zap.Error(err))
		return err
	}

	channelKey := request.MembershipKey.GroupKey.GetChannelKey()
	err = amqpChannel.ExchangeUnbind(ctx, membershipToDelete.UserKey.GetExchangeName(), "", amqp.WebsocketMessagesExchange, false, map[string]interface{}{
		"event_type": "chat.message",
		"channel_id": channelKey.String(),
		"x-match":    "all",
	})
	if err != nil {
		l.Error("could not unbind exchanges", zap.Error(err))
		return err
	}

	grp, err := g.groupStore.GetGroup(ctx, request.MembershipKey.GroupKey)
	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return err
	}

	text := fmt.Sprintf("%s has left #%s", usernameLeavingGroup, grp.Name)
	message := chat.NewContextBlock([]chat.BlockElement{
		chat.NewMarkdownObject(text)},
		nil,
	)

	_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(request.MembershipKey.GroupKey, membershipToDelete.UserKey, usernameLeavingGroup, text, []chat.Block{*message}, []chat.Attachment{}, nil))
	if err != nil {
		l.Error("could not send user leaving message", zap.Error(err))
		return err
	}

	return nil

}

func (g GroupService) RegisterUserAmqpSubscriptions(ctx context.Context) error {

	ctx, l := GetCtx(ctx, "GroupService", "RegisterUserAmqpSubscriptions")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session")
		return err
	}

	approvedMembershipStatus := group.ApprovedMembershipStatus
	memberships, err := g.groupStore.GetMembershipsForUser(ctx, userSession.GetUserKey(), &approvedMembershipStatus)
	if err != nil {
		l.Error("could not get active user memberships", zap.Error(err))
		return err
	}

	l.Debug(fmt.Sprintf("user has %d active subscriptions", len(memberships.Items)))

	for _, activeMembership := range memberships.Items {

		grp, err := g.groupStore.GetGroup(ctx, activeMembership.GetGroupKey())
		if err != nil {
			l.Error("could not get group", zap.Error(err))
			return err
		}

		channel, err := g.chatService.CreateChannel(ctx, grp.GetKey().GetChannelKey(), chat.GroupChannel)
		if err != nil {
			l.Error("could not create channel for group", zap.Error(err))
			return err
		}

		channelSubscriptionKey := model.NewChannelSubscriptionKey(channel.GetKey(), userSession.GetUserKey())
		_, err = g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
		if err != nil {
			l.Error("could not subscribe to channel", zap.Error(err))
			return err
		}

	}

	return nil
}
