package integration

import (
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	userreadmodels "github.com/commonpool/backend/pkg/auth/readmodel"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/eventstore"
	groupreadmodels "github.com/commonpool/backend/pkg/group/readmodels"
	resourcereadmodels "github.com/commonpool/backend/pkg/resource/readmodel"
	tradingreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	uuid "github.com/satori/go.uuid"
	_ "net/http/pprof"
	"sync"
)

var (
	userIncrementer   = 0
	groupCounter      = 0
	userIncrementerMu sync.Mutex
	createUserLock    sync.Mutex
)

func cleanDb() {

	session := srv.GraphDriver.GetSession()

	_, err := session.Run(`MATCH (n) DETACH DELETE n`, map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	srv.Db.Delete(chatstore.Channel{}, "1 = 1")
	srv.Db.Delete(chatstore.ChannelSubscription{}, "1 = 1")
	srv.Db.Delete(chatstore.Message{}, "1 = 1")
	srv.Db.Delete(eventstore.StreamEvent{}, "1 = 1")
	srv.Db.Delete(eventstore.Stream{}, "1 = 1")
	srv.Db.Delete(userreadmodels.UserReadModel{}, "1 = 1")
	srv.Db.Delete(groupreadmodels.MembershipReadModel{}, "1 = 1")
	srv.Db.Delete(groupreadmodels.GroupReadModel{}, "1 = 1")
	srv.Db.Delete(groupreadmodels.DBGroupUserReadModel{}, "1 = 1")
	srv.Db.Delete(resourcereadmodels.DbResourceReadModel{}, "1 = 1")
	srv.Db.Delete(resourcereadmodels.ResourceGroupNameReadModel{}, "1 = 1")
	srv.Db.Delete(resourcereadmodels.ResourceSharingReadModel{}, "1 = 1")
	srv.Db.Delete(tradingreadmodels.DBOfferReadModel{}, "1 = 1")
	srv.Db.Delete(tradingreadmodels.OfferItemReadModel{}, "1 = 1")
	srv.Db.Delete(tradingreadmodels.OfferResourceReadModel{}, "1 = 1")
	srv.Db.Delete(tradingreadmodels.OfferUserMembershipReadModel{}, "1 = 1")
	srv.Db.Delete(tradingreadmodels.OfferUserReadModel{}, "1 = 1")
}

func NewUser() *models.UserSession {
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
