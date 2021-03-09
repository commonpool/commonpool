package integration

import (
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	userreadmodels "github.com/commonpool/backend/pkg/auth/readmodel"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/eventstore"
	groupreadmodels "github.com/commonpool/backend/pkg/group/readmodels"
	resourcereadmodels "github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/commonpool/backend/pkg/server"
	tradingreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "net/http/pprof"
	"sync"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite
	server            *server.Server
	userIncrementer   int
	userIncrementerMu sync.Mutex
	createUserLock    sync.Mutex
	groupCounter      int
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.userIncrementer = 0
	s.userIncrementerMu = sync.Mutex{}
	s.groupCounter = 0
	srv, err := server.NewServer()
	if err != nil {
		print(err.Error())
	}
	s.server = srv
	s.cleanDb()
}

func (s *IntegrationTestSuite) NewUser() *models.UserSession {
	s.userIncrementerMu.Lock()
	defer func() {
		s.userIncrementerMu.Unlock()
	}()
	s.userIncrementer++
	var userId = uuid.NewV4().String()
	userEmail := fmt.Sprintf("user%d@email.com", s.userIncrementer)
	userName := fmt.Sprintf("user%d", s.userIncrementer)
	return &models.UserSession{
		Username:        userName,
		Subject:         userId,
		Email:           userEmail,
		IsAuthenticated: true,
	}
}

func (s *IntegrationTestSuite) cleanDb() {

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

func getDb(appConfig *config.AppConfig) *gorm.DB {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", appConfig.DbHost, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName, appConfig.DbPort)
	database, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return database
}
