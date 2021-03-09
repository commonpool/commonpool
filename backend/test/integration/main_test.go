package integration

import (
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	userreadmodels "github.com/commonpool/backend/pkg/auth/readmodel"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/eventstore"
	groupreadmodels "github.com/commonpool/backend/pkg/group/readmodels"
	resourcereadmodels "github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/commonpool/backend/pkg/server"
	tradingreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	_ "net/http/pprof"
	"sync"
	"testing"
)

var (
	userIncrementer   = 0
	groupCounter      = 0
	userIncrementerMu sync.Mutex
	createUserLock    sync.Mutex
)

type IntegrationTestBase struct {
	server  *server.Server
	servers []*server.Server
}

func (i *IntegrationTestBase) Setup() {
	srv, err := server.NewServer()
	if err != nil {
		print(err.Error())
	}
	i.server = srv

	for j := 0; j < 10; j++ {
		srv, err = server.NewServer()
		if err != nil {
			panic(err)
		}
		i.servers = append(i.servers, srv)
	}

	i.cleanDb()
}

func (s *IntegrationTestBase) cleanDb() {

	session := s.server.GraphDriver.GetSession()

	_, err := session.Run(`MATCH (n) DETACH DELETE n`, map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	s.server.Db.Delete(chatstore.Channel{}, "1 = 1")
	s.server.Db.Delete(chatstore.ChannelSubscription{}, "1 = 1")
	s.server.Db.Delete(chatstore.Message{}, "1 = 1")
	s.server.Db.Delete(eventstore.StreamEvent{}, "1 = 1")
	s.server.Db.Delete(eventstore.Stream{}, "1 = 1")
	s.server.Db.Delete(userreadmodels.UserReadModel{}, "1 = 1")
	s.server.Db.Delete(groupreadmodels.MembershipReadModel{}, "1 = 1")
	s.server.Db.Delete(groupreadmodels.GroupReadModel{}, "1 = 1")
	s.server.Db.Delete(groupreadmodels.DBGroupUserReadModel{}, "1 = 1")
	s.server.Db.Delete(resourcereadmodels.DbResourceReadModel{}, "1 = 1")
	s.server.Db.Delete(resourcereadmodels.ResourceGroupNameReadModel{}, "1 = 1")
	s.server.Db.Delete(resourcereadmodels.ResourceSharingReadModel{}, "1 = 1")
	s.server.Db.Delete(tradingreadmodels.DBOfferReadModel{}, "1 = 1")
	s.server.Db.Delete(tradingreadmodels.OfferItemReadModel{}, "1 = 1")
	s.server.Db.Delete(tradingreadmodels.OfferResourceReadModel{}, "1 = 1")
	s.server.Db.Delete(tradingreadmodels.OfferUserMembershipReadModel{}, "1 = 1")
	s.server.Db.Delete(tradingreadmodels.OfferUserReadModel{}, "1 = 1")
}

type IntegrationTestSuite struct {
	suite.Suite
	*IntegrationTestBase
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.IntegrationTestBase = &IntegrationTestBase{}
	s.IntegrationTestBase.Setup()
}

func (s *IntegrationTestBase) NewUser() *models.UserSession {
	userIncrementerMu.Lock()
	defer func() {
		userIncrementerMu.Unlock()
	}()
	userIncrementer++
	var userId = uuid.NewV4().String()
	userEmail := fmt.Sprintf("user%d@email.com", userIncrementer)
	userName := fmt.Sprintf("user%d", userIncrementer)
	return &models.UserSession{
		Username:        userName,
		Subject:         userId,
		Email:           userEmail,
		IsAuthenticated: true,
	}
}
