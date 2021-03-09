package listeners

import (
	"context"
	"encoding/json"
	authdomain "github.com/commonpool/backend/pkg/auth/domain"
	db2 "github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
	groupdomain "github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	resourcedomain "github.com/commonpool/backend/pkg/resource/domain"
	tradingdomain "github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
	readmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type Aggregate interface {
	keys.StreamKeyer
	eventsource.ChangeGetter
	eventsource.RevisionGetter
	MarkAsCommitted()
}

type OfferReadModelTestSuite struct {
	suite.Suite
	db                   *gorm.DB
	eventStore           eventstore.EventStore
	listener             *OfferReadModelHandler
	getOffer             *queries.GetOffer
	getOffers            *queries.GetOffers
	getOfferItem         *queries.GetOfferItem
	getOfferItemOfferKey *queries.GetOfferKeyForOfferItemKey
}

func (s *OfferReadModelTestSuite) SetupSuite() {
	db := db2.NewTestDb("OfferReadModelTestSuite")
	s.db = db
	listener := &OfferReadModelHandler{
		db: db,
	}
	if err := listener.migrateDatabase(); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	s.listener = listener

	eventMapper := eventsource.NewEventMapper()
	if err := groupdomain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := authdomain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := resourcedomain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := tradingdomain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}

	eventStore := postgres.NewPostgresEventStore(db, eventMapper)
	if err := eventStore.MigrateDatabase(); err != nil {
		s.FailNow(err.Error())
	}
	s.eventStore = eventStore

	// s.cleanDb(db)

	s.getOffer = queries.NewGetOffer(db)
	s.getOffers = queries.NewGetOffers(db)
	s.getOfferItem = queries.NewGetOfferItem(db)
	s.getOfferItemOfferKey = queries.NewGetOfferKeyForOfferItemKey(db)
}

func (s *OfferReadModelTestSuite) cleanDb(db *gorm.DB) {
	if err := db.Delete(&readmodels.OfferGroupReadModel{}, "1 = 1").Error; !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := db.Delete(&readmodels.OfferItemReadModel{}, "1 = 1").Error; !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := db.Delete(&readmodels.OfferUserReadModel{}, "1 = 1").Error; !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := db.Delete(&readmodels.OfferResourceReadModel{}, "1 = 1").Error; !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := db.Delete(&readmodels.DBOfferReadModel{}, "1 = 1").Error; !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := db.Delete(&eventstore.StreamEvent{}, "1 = 1").Error; !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := db.Delete(&eventstore.Stream{}, "1 = 1").Error; !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
}

func (s *OfferReadModelTestSuite) saveAndApply(aggregate Aggregate) ([]eventsource.Event, error) {
	evts, err := s.eventStore.Save(context.TODO(), aggregate.StreamKey(), aggregate.GetVersion(), aggregate.GetChanges())
	if err != nil {
		return nil, err
	}
	err = s.listener.HandleEvents(context.TODO(), evts)
	aggregate.MarkAsCommitted()
	return evts, err
}

func (s *OfferReadModelTestSuite) discoverUser(user *authdomain.User, userInfo authdomain.UserInfo) ([]eventsource.Event, error) {
	err := user.DiscoverUser(userInfo)
	if !assert.NoError(s.T(), err) {
		return nil, nil
	}
	return s.saveAndApply(user)
}

func (s *OfferReadModelTestSuite) createGroup(group *groupdomain.Group, createdBy keys.UserKey, groupInfo groupdomain.GroupInfo) ([]eventsource.Event, error) {
	err := group.CreateGroup(createdBy, groupInfo)
	if !assert.NoError(s.T(), err) {
		return nil, nil
	}
	return s.saveAndApply(group)
}

func (s *OfferReadModelTestSuite) registerResource(
	resource *resourcedomain.Resource,
	createdBy keys.UserKey,
	target keys.Targetter,
	resourceInfo resourcedomain.ResourceInfo,
	groupKeys ...keys.GroupKey) ([]eventsource.Event, error) {
	err := resource.Register(createdBy, target.Target(), resourceInfo, keys.NewGroupKeys(groupKeys))
	if !assert.NoError(s.T(), err) {
		return nil, nil
	}
	return s.saveAndApply(resource)
}

func (s *OfferReadModelTestSuite) Test() {

	var err error

	user1Key := keys.GenerateUserKey()
	user1Info := authdomain.NewUserInfo().WithUsername("Test1").WithEmail("Test1@example.com")
	user1 := authdomain.NewUser(user1Key)
	if _, err := s.discoverUser(user1, user1Info); !assert.NoError(s.T(), err) {
		return
	}

	user2Key := keys.GenerateUserKey()
	user2Info := authdomain.NewUserInfo().WithUsername("Test2").WithEmail("Test2@example.com")
	user2 := authdomain.NewUser(user2Key)
	if _, err := s.discoverUser(user2, user2Info); !assert.NoError(s.T(), err) {
		return
	}

	group1Key := keys.GenerateGroupKey()
	group1 := groupdomain.NewGroup(group1Key)
	group1Info := groupdomain.NewGroupInfo().WithName("Test Group 1").WithDescription("Test Group 1 Description")
	if _, err := s.createGroup(group1, user1Key, group1Info); !assert.NoError(s.T(), err) {
		return
	}

	group2Key := keys.GenerateGroupKey()
	group2 := groupdomain.NewGroup(group2Key)
	group2Info := groupdomain.NewGroupInfo().
		WithName("Test Group 2").
		WithDescription("Test Group 2Description")
	if _, err := s.createGroup(group2, user2Key, group2Info); !assert.NoError(s.T(), err) {
		return
	}

	resource1Key := keys.GenerateResourceKey()
	resource1Info := resourcedomain.NewResourceInfo().
		WithIsOffer().
		WithIsService().
		WithName("Test Resource 1").
		WithDescription("Test Resource 1 Description").
		WithValue(resourcedomain.NewResourceValueEstimation().
			WithFromToValueType().WithHoursFromTo(1, 2))
	resource1 := resourcedomain.NewResource(resource1Key)
	if _, err := s.registerResource(resource1, user1Key, user1Key, resource1Info); !assert.NoError(s.T(), err) {
		return
	}

	resource2Key := keys.GenerateResourceKey()
	resource2Info := resourcedomain.NewResourceInfo().
		WithIsOffer().
		WithIsService().
		WithName("Test Resource 2").
		WithDescription("Test Resource 2 Description").
		WithValue(resourcedomain.NewResourceValueEstimation().
			WithFromToValueType().WithHoursFromTo(1, 2))
	resource2 := resourcedomain.NewResource(resource2Key)
	if _, err := s.registerResource(resource2, user1Key, user2Key, resource2Info); !assert.NoError(s.T(), err) {
		return
	}

	resource3Key := keys.GenerateResourceKey()
	resource3Info := resourcedomain.NewResourceInfo().
		WithIsOffer().
		WithIsService().
		WithName("Test Resource 3").
		WithDescription("Test Resource 3 Description").
		WithValue(resourcedomain.NewResourceValueEstimation().
			WithFromToValueType().WithHoursFromTo(2, 3))
	resource3 := resourcedomain.NewResource(resource3Key)
	if _, err := s.registerResource(resource3, user2Key, user2Key, resource3Info); !assert.NoError(s.T(), err) {
		return
	}

	offerKey := keys.GenerateOfferKey()
	offer := tradingdomain.NewOffer(offerKey)
	offerItemKey := keys.GenerateOfferItemKey()
	offerItem2Key := keys.GenerateOfferItemKey()
	offerItem3Key := keys.GenerateOfferItemKey()
	offerItem4Key := keys.GenerateOfferItemKey()
	err = offer.Submit(user1Key, group1Key, tradingdomain.NewSubmitOfferItems(
		tradingdomain.NewResourceTransferItemInput(
			offerItemKey,
			group1Key,
			resource1Key,
		),
		tradingdomain.NewProvideServiceItemInput(
			offerItem2Key,
			user1Key,
			user2Key,
			resource2Key,
			2*time.Hour,
		),
		tradingdomain.NewBorrowResourceInput(
			offerItem3Key,
			user2Key,
			resource3Key,
			2*time.Hour,
		),
		tradingdomain.NewCreditTransferItemInput(
			offerItem4Key,
			group1Key,
			user1Key,
			time.Hour*5,
		),
	),
	)
	if !assert.NoError(s.T(), err) {
		return
	}
	evts, err := s.saveAndApply(offer)
	if !assert.NoError(s.T(), err) {
		return
	}

	user1Info = user1Info.WithUsername("Test1-2")
	err = user1.ChangeUserInfo(user1Info)
	if !assert.NoError(s.T(), err) {
		return
	}
	_, err = s.saveAndApply(user1)
	if !assert.NoError(s.T(), err) {
		return
	}

	resource1Info = resource1Info.WithName("Changed Resource Name")
	if err := resource1.ChangeInfo(user2Key, resource1Info.AsUpdate()); !assert.NoError(s.T(), err) {
		return
	}
	_, err = s.saveAndApply(resource1)
	if !assert.NoError(s.T(), err) {
		return
	}

	group1Info = group1Info.WithName("Changed Group Name")
	if err := group1.ChangeInfo(user2Key, group1Info); !assert.NoError(s.T(), err) {
		return
	}
	_, err = s.saveAndApply(group1)
	if !assert.NoError(s.T(), err) {
		return
	}

	if err := offer.ApproveOfferItem(user1Key, offerItemKey, tradingdomain.Inbound, tradingdomain.ApproveAllMatrix); !assert.NoError(s.T(), err) {
		return
	}
	approval1Events, err := s.saveAndApply(offer)
	if !assert.NoError(s.T(), err) {
		return
	}

	if err := offer.ApproveOfferItem(user2Key, offerItemKey, tradingdomain.Outbound, tradingdomain.ApproveAllMatrix); !assert.NoError(s.T(), err) {
		return
	}
	approval2Events, err := s.saveAndApply(offer)
	if !assert.NoError(s.T(), err) {
		return
	}

	rm, err := s.getOffer.Get(context.Background(), offerKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	all, err := s.getOffers.Get(context.Background(), user1Key)
	if !assert.NoError(s.T(), err) {
		return
	}
	js, _ := json.MarshalIndent(all, "", " ")
	s.T().Log("\n" + string(js))

	all, err = s.getOffers.Get(context.Background(), user2Key)
	if !assert.NoError(s.T(), err) {
		return
	}
	js, _ = json.MarshalIndent(all, "", " ")
	s.T().Log("\n" + string(js))

	intPtr := func(i int) *int {
		return &i
	}
	strPtr := func(s string) *string {
		return &s
	}
	durationPtr := func(d time.Duration) *time.Duration {
		return &d
	}
	timePtr := func(d time.Time) *time.Time {
		return &d
	}
	expected := readmodels.OfferReadModel{
		OfferReadModelBase: readmodels.OfferReadModelBase{
			OfferKey:    offerKey,
			GroupKey:    group1Key,
			Status:      tradingdomain.Pending,
			Version:     2,
			DeclinedAt:  nil,
			SubmittedAt: evts[0].GetEventTime(),
			ApprovedAt:  nil,
			CompletedAt: nil,
		},
		DeclinedBy: nil,
		SubmittedBy: &readmodels.OfferUserReadModel{
			UserKey:  user1Key,
			Username: user1Info.Username,
			Version:  1,
		},
		OfferItems: []*readmodels.OfferItemReadModel2{
			{
				OfferItemReadModelBase: readmodels.OfferItemReadModelBase{
					OfferItemKey:       offerItemKey,
					OfferKey:           offerKey,
					Type:               tradingdomain.ResourceTransfer,
					ApprovedInbound:    true,
					ApprovedInboundAt:  timePtr(approval1Events[0].GetEventTime()),
					ApprovedOutbound:   true,
					ApprovedOutboundAt: timePtr(approval2Events[0].GetEventTime()),
					Version:            2,
				},
				To: &readmodels.OfferItemTargetReadModel{
					Target: keys.Target{
						GroupKey: &group1Key,
						Type:     keys.GroupTarget,
					},
					GroupName:    strPtr("Changed Group Name"),
					GroupVersion: intPtr(2),
				},
				Resource: &readmodels.OfferResourceReadModel{
					ResourceKey:  resource1Key,
					ResourceName: "Changed Resource Name",
					Version:      1,
					ResourceType: resourcedomain.ServiceResource,
					CallType:     resourcedomain.Offer,
					Owner: keys.Target{
						UserKey: &user1Key,
						Type:    keys.UserTarget,
					},
				},
				ApprovedInboundBy: &readmodels.OfferUserReadModel{
					UserKey:  user1Key,
					Username: "Test1-2",
					Version:  1,
				},
				ApprovedOutboundBy: &readmodels.OfferUserReadModel{
					UserKey:  user2Key,
					Username: "Test2",
					Version:  0,
				},
			}, {
				OfferItemReadModelBase: readmodels.OfferItemReadModelBase{
					OfferItemKey: offerItem2Key,
					OfferKey:     offerKey,
					Type:         tradingdomain.ProvideService,
					Duration:     durationPtr(2 * time.Hour),
				},
				To: &readmodels.OfferItemTargetReadModel{
					Target: keys.Target{
						UserKey: &user2Key,
						Type:    keys.UserTarget,
					},
					UserName:    strPtr("Test2"),
					UserVersion: intPtr(0),
				},
				Resource: &readmodels.OfferResourceReadModel{
					ResourceKey:  resource2Key,
					ResourceName: "Test Resource 2",
					Version:      0,
					ResourceType: resourcedomain.ServiceResource,
					CallType:     resourcedomain.Offer,
					Owner: keys.Target{
						UserKey: &user2Key,
						Type:    keys.UserTarget,
					},
				},
			}, {
				OfferItemReadModelBase: readmodels.OfferItemReadModelBase{
					OfferItemKey: offerItem3Key,
					OfferKey:     offerKey,
					Type:         tradingdomain.BorrowResource,
					Duration:     durationPtr(2 * time.Hour),
				},
				To: &readmodels.OfferItemTargetReadModel{
					Target: keys.Target{
						UserKey: &user2Key,
						Type:    keys.UserTarget,
					},
					UserName:    strPtr("Test2"),
					UserVersion: intPtr(0),
				},
				Resource: &readmodels.OfferResourceReadModel{
					ResourceKey:  resource3Key,
					ResourceName: "Test Resource 3",
					Version:      0,
					ResourceType: resourcedomain.ServiceResource,
					CallType:     resourcedomain.Offer,
					Owner: keys.Target{
						UserKey: &user2Key,
						Type:    keys.UserTarget,
					},
				},
			}, {
				OfferItemReadModelBase: readmodels.OfferItemReadModelBase{
					OfferItemKey: offerItem4Key,
					OfferKey:     offerKey,
					Type:         tradingdomain.CreditTransfer,
					Amount:       durationPtr(5 * time.Hour),
				},
				To: &readmodels.OfferItemTargetReadModel{
					Target: keys.Target{
						UserKey: &user1Key,
						Type:    keys.UserTarget,
					},
					UserName:    strPtr("Test1-2"),
					UserVersion: intPtr(1),
				},
				From: &readmodels.OfferItemTargetReadModel{
					Target: keys.Target{
						GroupKey: &group1Key,
						Type:     keys.GroupTarget,
					},
					GroupName:    strPtr("Changed Group Name"),
					GroupVersion: intPtr(2),
				},
			},
		},
	}

	expectedJson, err := json.MarshalIndent(expected, "", "  ")
	if !assert.NoError(s.T(), err) {
		return
	}
	s.T().Log(string(expectedJson))

	actualJson, err := json.MarshalIndent(rm, "", "  ")
	if !assert.NoError(s.T(), err) {
		return
	}
	s.T().Log(string(actualJson))

	assert.Equal(s.T(), string(expectedJson), string(actualJson))

}

func TestReadModelSuite(t *testing.T) {
	suite.Run(t, &OfferReadModelTestSuite{})
}
