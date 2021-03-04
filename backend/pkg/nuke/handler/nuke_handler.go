package handler

import (
	"github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/graph"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/mq"
	store2 "github.com/commonpool/backend/pkg/transaction/store"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	amqpClient  mq.Client
	graphClient graph.Driver
}

func NewHandler(db *gorm.DB, amqpClient mq.Client, graphClient graph.Driver) *Handler {
	return &Handler{
		db:          db,
		amqpClient:  amqpClient,
		graphClient: graphClient,
	}
}

func (h *Handler) Nuke(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "Nuke")

	var channelSubscriptions []*store.ChannelSubscription

	if err := h.db.Where("1 = 1").Find(&channelSubscriptions).Error; err != nil {
		return err
	}
	var userKeys = keys.NewEmptyUserKeys()
	for _, subscription := range channelSubscriptions {
		userKeys.Append(subscription.Map().UserKey)
	}

	ch, err := h.amqpClient.GetChannel()
	if err != nil {
		return err
	}
	for _, userKey := range userKeys.Items {
		err := ch.ExchangeDelete(ctx, userKey.GetExchangeName(), false, false)
		if err != nil {
			return err
		}
	}

	if err := ch.ExchangeDelete(ctx, mq.MessagesExchange, false, false); err != nil {
		return err
	}

	if err := ch.ExchangeDelete(ctx, mq.WebsocketMessagesExchange, false, false); err != nil {
		return err
	}

	if err := h.db.Where("1 = 1").Delete(&store.Channel{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&store.ChannelSubscription{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&store.Message{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&store2.TransactionEntry{}).Error; err != nil {
		return err
	}

	dbSession := h.graphClient.GetSession()
	res, err := dbSession.Run(`MATCH (n) DETACH DELETE n`, map[string]interface{}{})
	if err != nil {
		return err
	}

	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (h *Handler) Register(g *echo.Group) {
	g.GET("/nuke", h.Nuke)
}
