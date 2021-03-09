package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/test"

	res "github.com/commonpool/backend/pkg/resource/handler"
	trd "github.com/commonpool/backend/pkg/trading/handler"
	"time"
)

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenUsers() {

	ctx := context.Background()

	user1, user1Cli := s.testUserCli(s.T())
	user2, _ := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !s.NoError(s.testGroup2(s.T(), user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !s.NoError(user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !s.NoError(user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	s.Equal(1, len(offerResponse.Offer.OfferItems))

}

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenUsersAndGroup() {

	ctx := context.Background()

	user1, user1Cli := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !s.NoError(s.testGroup2(s.T(), user1, &group)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !s.NoError(user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !s.NoError(user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, group.Group, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	s.Equal(1, len(offerResponse.Offer.OfferItems))

}

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenGroupAndMultipleUsers() {

	ctx := context.Background()

	user1, user1Cli := s.testUserCli(s.T())
	user2, _ := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !s.NoError(s.testGroup2(s.T(), user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !s.NoError(user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !s.NoError(user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
		domain.NewCreditTransferItemInputBase(user2, group.Group, time.Hour*2),
		domain.NewCreditTransferItemInputBase(group.Group, user1, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	s.Equal(3, len(offerResponse.Offer.OfferItems))

}

func (s *IntegrationTestSuite) TestUsersCanAcceptOfferBetweenUsers() {

	ctx := context.Background()

	user1, user1Cli := s.testUserCli(s.T())
	user2, user2Cli := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !s.NoError(s.testGroup2(s.T(), user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !s.NoError(user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offer := &trd.GetOfferResponse{}
	if !s.NoError(user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
		domain.NewCreditTransferItemInputBase(group, user1, time.Hour*2),
	).AsRequest(), offer)) {
		return
	}
	if !s.Equal(2, len(offer.Offer.OfferItems)) {
		return
	}

	time.Sleep(1 * time.Second)

	if !s.NoError(user1Cli.AcceptOffer(ctx, offer)) {
		return
	}

	time.Sleep(1 * time.Second)

	if !s.NoError(user2Cli.AcceptOffer(ctx, offer)) {
		return
	}

}

func (s *IntegrationTestSuite) TestUserCannotCreateOfferForResourceNotSharedWithGroup() {

	ctx := context.Background()
	user1, user1Cli := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !s.NoError(s.testGroup2(s.T(), user1, &group)) {
		return
	}

	var resource res.GetResourceResponse
	if !s.NoError(user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resource)) {
		return
	}

	var offer trd.GetOfferResponse
	err := user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewResourceTransferItemInputBase(group, resource)).AsRequest(), &offer)
	if !s.Error(err) {
		return
	}
	if !s.ErrorIs(exceptions.ErrForbidden, err) {
		return
	}

}

func (s *IntegrationTestSuite) TestCannotCreateResourceTransferItemForResourceAlreadyOwned() {

	ctx := context.Background()
	user, cli := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !s.NoError(s.testGroup2(s.T(), user, &group)) {
		return
	}

	var resource res.GetResourceResponse
	if !s.NoError(s.testResource(ctx, cli, &resource, group)) {
		return
	}

	var offer trd.GetOfferResponse
	err := cli.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewResourceTransferItemInputBase(user, resource)).AsRequest(), &offer)

	if !s.Error(err) {
		return
	}

	if !s.ErrorIs(exceptions.ErrBadRequest("OfferItem Resource destination is the same as the resource owner"), err) {
		return
	}
}

func (s *IntegrationTestSuite) TestUsersCanDeclineOffer() {

	ctx := context.Background()

	user1, cli1 := s.testUserCli(s.T())
	user2, cli2 := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !s.NoError(s.testGroup2(s.T(), user1, &group)) {
		return
	}

	var offer trd.GetOfferResponse
	if !s.NoError(cli1.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2)).AsRequest(), &offer)) {
		return
	}

	time.Sleep(100 * time.Millisecond)

	if !s.NoError(cli1.AcceptOffer(ctx, offer)) {
		return
	}

	time.Sleep(100 * time.Millisecond)

	if !s.NoError(cli1.GetOffer(ctx, offer, &offer)) {
		return
	}

	s.Equal(domain.Pending, offer.Offer.Status)

	time.Sleep(100 * time.Millisecond)

	if !s.NoError(cli2.DeclineOffer(ctx, offer)) {
		return
	}

	time.Sleep(100 * time.Millisecond)

	if !s.NoError(cli1.GetOffer(ctx, offer, &offer)) {
		return
	}

	s.Equal(domain.Declined, offer.Offer.Status)

}
