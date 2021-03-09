package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/test"

	res "github.com/commonpool/backend/pkg/resource/handler"
	trd "github.com/commonpool/backend/pkg/trading/handler"
	"github.com/stretchr/testify/assert"
	"time"
)

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenUsers() {

	ctx := context.Background()

	user1, user1Cli := s.testUserCli(s.T())
	user2, _ := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !assert.NoError(s.T(), s.testGroup2(s.T(), user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !assert.NoError(s.T(), user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !assert.NoError(s.T(), user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	assert.Equal(s.T(), 1, len(offerResponse.Offer.OfferItems))

}

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenUsersAndGroup() {

	ctx := context.Background()

	user1, user1Cli := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !assert.NoError(s.T(), s.testGroup2(s.T(), user1, &group)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !assert.NoError(s.T(), user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !assert.NoError(s.T(), user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, group.Group, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	assert.Equal(s.T(), 1, len(offerResponse.Offer.OfferItems))

}

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenGroupAndMultipleUsers() {

	ctx := context.Background()

	user1, user1Cli := s.testUserCli(s.T())
	user2, _ := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !assert.NoError(s.T(), s.testGroup2(s.T(), user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !assert.NoError(s.T(), user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !assert.NoError(s.T(), user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
		domain.NewCreditTransferItemInputBase(user2, group.Group, time.Hour*2),
		domain.NewCreditTransferItemInputBase(group.Group, user1, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	assert.Equal(s.T(), 3, len(offerResponse.Offer.OfferItems))

}

func (s *IntegrationTestSuite) TestUsersCanAcceptOfferBetweenUsers() {

	ctx := context.Background()

	user1, user1Cli := s.testUserCli(s.T())
	user2, user2Cli := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !assert.NoError(s.T(), s.testGroup2(s.T(), user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !assert.NoError(s.T(), user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offer := &trd.GetOfferResponse{}
	if !assert.NoError(s.T(), user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
		domain.NewCreditTransferItemInputBase(group, user1, time.Hour*2),
	).AsRequest(), offer)) {
		return
	}
	if !assert.Equal(s.T(), 2, len(offer.Offer.OfferItems)) {
		return
	}
	if !assert.NoError(s.T(), user1Cli.AcceptOffer(ctx, offer)) {
		return
	}
	if !assert.NoError(s.T(), user2Cli.AcceptOffer(ctx, offer)) {
		return
	}

}

func (s *IntegrationTestSuite) TestUserCannotCreateOfferForResourceNotSharedWithGroup() {
	s.T().Parallel()

	ctx := context.Background()
	user1, user1Cli := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !assert.NoError(s.T(), s.testGroup2(s.T(), user1, &group)) {
		return
	}

	var resource res.GetResourceResponse
	if !assert.NoError(s.T(), user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resource)) {
		return
	}

	var offer trd.GetOfferResponse
	err := user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewResourceTransferItemInputBase(group, resource)).AsRequest(), &offer)
	if !assert.Error(s.T(), err) {
		return
	}
	if !assert.ErrorIs(s.T(), exceptions.ErrForbidden, err) {
		return
	}

}

func (s *IntegrationTestSuite) TestCannotCreateResourceTransferItemForResourceAlreadyOwned() {
	s.T().Parallel()

	ctx := context.Background()
	user, cli := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !assert.NoError(s.T(), s.testGroup2(s.T(), user, &group)) {
		return
	}

	var resource res.GetResourceResponse
	if !assert.NoError(s.T(), s.testResource(ctx, cli, &resource, group)) {
		return
	}

	var offer trd.GetOfferResponse
	err := cli.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewResourceTransferItemInputBase(user, resource)).AsRequest(), &offer)

	if !assert.Error(s.T(), err) {
		return
	}

	if !assert.ErrorIs(s.T(), exceptions.ErrForbidden, err) {
		return
	}
}

func (s *IntegrationTestSuite) TestUsersCanDeclineOffer() {
	s.T().Parallel()
	ctx := context.Background()

	user1, cli1 := s.testUserCli(s.T())
	user2, cli2 := s.testUserCli(s.T())

	var group grouphandler.GetGroupResponse
	if !assert.NoError(s.T(), s.testGroup2(s.T(), user1, &group)) {
		return
	}

	var offer trd.GetOfferResponse
	if !assert.NoError(s.T(), cli1.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2)).AsRequest(), &offer)) {
		return
	}

	if !assert.NoError(s.T(), cli1.AcceptOffer(ctx, offer)) {
		return
	}

	if !assert.NoError(s.T(), cli1.GetOffer(ctx, offer, &offer)) {
		return
	}

	assert.Equal(s.T(), domain.Pending, offer.Offer.Status)

	if !assert.NoError(s.T(), cli2.DeclineOffer(ctx, offer)) {
		return
	}

	assert.Equal(s.T(), domain.Declined, offer.Offer.Status)

}
