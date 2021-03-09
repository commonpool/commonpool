package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/test"
	"github.com/stretchr/testify/assert"
	"testing"

	res "github.com/commonpool/backend/pkg/resource/handler"
	trd "github.com/commonpool/backend/pkg/trading/handler"
	"time"
)

func TestUserCanSubmitOfferBetweenUsers(t *testing.T) {

	ctx := context.Background()

	user1, user1Cli := testUserCli(t)
	user2, _ := testUserCli(t)

	var group grouphandler.GetGroupResponse
	if !assert.NoError(t, testGroup2(t, user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !assert.NoError(t, user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !assert.NoError(t, user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	assert.Equal(t, 1, len(offerResponse.Offer.OfferItems))

}

func TestUserCanSubmitOfferBetweenUsersAndGroup(t *testing.T) {

	ctx := context.Background()

	user1, user1Cli := testUserCli(t)

	var group grouphandler.GetGroupResponse
	if !assert.NoError(t, testGroup2(t, user1, &group)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !assert.NoError(t, user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !assert.NoError(t, user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, group.Group, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	assert.Equal(t, 1, len(offerResponse.Offer.OfferItems))

}

func TestUserCanSubmitOfferBetweenGroupAndMultipleUsers(t *testing.T) {

	ctx := context.Background()

	user1, user1Cli := testUserCli(t)
	user2, _ := testUserCli(t)

	var group grouphandler.GetGroupResponse
	if !assert.NoError(t, testGroup2(t, user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !assert.NoError(t, user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offerResponse := &trd.GetOfferResponse{}
	if !assert.NoError(t, user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
		domain.NewCreditTransferItemInputBase(user2, group.Group, time.Hour*2),
		domain.NewCreditTransferItemInputBase(group.Group, user1, time.Hour*2),
	).AsRequest(), offerResponse)) {
		return
	}

	assert.Equal(t, 3, len(offerResponse.Offer.OfferItems))

}

func TestUsersCanAcceptOfferBetweenUsers(t *testing.T) {

	ctx := context.Background()

	user1, user1Cli := testUserCli(t)
	user2, user2Cli := testUserCli(t)

	var group grouphandler.GetGroupResponse
	if !assert.NoError(t, testGroup2(t, user1, &group, user2)) {
		return
	}

	var resourceResponse res.GetResourceResponse
	if !assert.NoError(t, user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resourceResponse)) {
		return
	}

	time.Sleep(1 * time.Second)

	offer := &trd.GetOfferResponse{}
	if !assert.NoError(t, user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(
		group.Group.GroupKey,
		domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2),
		domain.NewCreditTransferItemInputBase(group, user1, time.Hour*2),
	).AsRequest(), offer)) {
		return
	}
	if !assert.Equal(t, 2, len(offer.Offer.OfferItems)) {
		return
	}

	time.Sleep(1 * time.Second)

	if !assert.NoError(t, user1Cli.AcceptOffer(ctx, offer)) {
		return
	}

	time.Sleep(1 * time.Second)

	if !assert.NoError(t, user2Cli.AcceptOffer(ctx, offer)) {
		return
	}

}

func TestUserCannotCreateOfferForResourceNotSharedWithGroup(t *testing.T) {

	ctx := context.Background()
	user1, user1Cli := testUserCli(t)

	var group grouphandler.GetGroupResponse
	if !assert.NoError(t, testGroup2(t, user1, &group)) {
		return
	}

	var resource res.GetResourceResponse
	if !assert.NoError(t, user1Cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resource)) {
		return
	}

	time.Sleep(1 * time.Second)

	var offer trd.GetOfferResponse
	err := user1Cli.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewResourceTransferItemInputBase(group, resource)).AsRequest(), &offer)
	if !assert.Error(t, err) {
		return
	}
	if !assert.ErrorIs(t, exceptions.ErrBadRequest(""), err) {
		return
	}

}

func TestCannotCreateResourceTransferItemForResourceAlreadyOwned(t *testing.T) {

	ctx := context.Background()
	user, cli := testUserCli(t)

	var group grouphandler.GetGroupResponse
	if !assert.NoError(t, testGroup2(t, user, &group)) {
		return
	}

	var resource res.GetResourceResponse
	if !assert.NoError(t, testResource(ctx, cli, &resource, group)) {
		return
	}

	var offer trd.GetOfferResponse
	err := cli.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewResourceTransferItemInputBase(user, resource)).AsRequest(), &offer)

	if !assert.Error(t, err) {
		return
	}

	if !assert.ErrorIs(t, exceptions.ErrBadRequest("OfferItem Resource destination is the same as the resource owner"), err) {
		return
	}
}

func TestUsersCanDeclineOffer(t *testing.T) {

	ctx := context.Background()

	user1, cli1 := testUserCli(t)
	user2, cli2 := testUserCli(t)

	var group grouphandler.GetGroupResponse
	if !assert.NoError(t, testGroup2(t, user1, &group)) {
		return
	}

	var offer trd.GetOfferResponse
	if !assert.NoError(t, cli1.SubmitOffer(ctx, trd.NewSendOfferPayload(group, domain.NewCreditTransferItemInputBase(user1, user2, time.Hour*2)).AsRequest(), &offer)) {
		return
	}

	time.Sleep(100 * time.Millisecond)

	if !assert.NoError(t, cli1.AcceptOffer(ctx, offer)) {
		return
	}

	time.Sleep(100 * time.Millisecond)

	if !assert.NoError(t, cli1.GetOffer(ctx, offer, &offer)) {
		return
	}

	assert.Equal(t, domain.Pending, offer.Offer.Status)

	time.Sleep(100 * time.Millisecond)

	if !assert.NoError(t, cli2.DeclineOffer(ctx, offer)) {
		return
	}

	time.Sleep(100 * time.Millisecond)

	if !assert.NoError(t, cli1.GetOffer(ctx, offer, &offer)) {
		return
	}

	assert.Equal(t, domain.Declined, offer.Offer.Status)

}
