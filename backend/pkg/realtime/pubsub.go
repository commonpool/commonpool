package realtime

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/commonpool/backend/pkg/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	amqp          mq.Client
	chatService   chat.Service
	authorization auth.Authenticator
}

func NewRealtimeHandler(amqpClient mq.Client, chatService chat.Service, authorization auth.Authenticator) *Handler {
	return &Handler{
		amqp:          amqpClient,
		chatService:   chatService,
		authorization: authorization,
	}
}

func (h *Handler) Register(g *echo.Group) {
	g.GET("/ws", h.Websocket, h.authorization.Authenticate(false))
}

func (h *Handler) websocketAnonymous(ctx context.Context, response *echo.Response, request *http.Request) error {

	ctx = logging.NewContext(ctx)
	l := logging.WithContext(ctx)

	l.Debug("Anonymous websocket")

	ws, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		l.Error("could not upgrade websocket connection", zap.Error(err))
		return err
	}
	defer ws.Close()

	hub := newHub()
	go hub.run()

	client := NewAnonymousClient(hub, ws)

	go client.writePump()
	go client.readPump()

	redirectResponse, err := h.authorization.GetRedirectResponse(request)
	if err != nil {
		return err
	}

	jsBytes, err := json.Marshal(redirectResponse)
	if err != nil {
		return err
	}

	client.send <- jsBytes

	time.Sleep(time.Minute * 10)

	return nil

}

func (h *Handler) Websocket(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "Websocket")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	userSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return h.websocketAnonymous(ctx, c.Response(), c.Request())
	}
	userKey := userSession.GetUserKey()

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		l.Error("could not upgrade websocket connection", zap.Error(err))
		return err
	}
	defer ws.Close()

	hub := newHub()
	go hub.run()

	consumerKey := utils.ShortUuid(uuid.NewV4())
	queueName := "chat.ws." + userKey.String() + "." + consumerKey

	amqpChannel, err := h.amqp.GetChannel()
	if err != nil {
		l.Error("cold not get amqp channel", zap.Error(err))
		return err
	}
	defer amqpChannel.Close()

	userExchangeName, err := h.chatService.CreateUserExchange(ctx, userKey)
	if err != nil {
		l.Error("could not create user exchange", zap.Error(err))
		return err
	}

	err = amqpChannel.QueueDeclare(ctx, queueName, false, true, false, false, nil)
	if err != nil {
		l.Error("could not declare websocket amqp queue", zap.Error(err))
		return err
	}

	err = amqpChannel.QueueBind(ctx, queueName, "", userExchangeName, false, nil)
	if err != nil {
		l.Error("could not bind consumer queue to exchange", zap.Error(err))
		return err
	}

	client := NewClient(hub, ws, amqpChannel, &queueName, &consumerKey)

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	return client.eventPump(ctx)

	// ch, err := h.mq.RegisterWsChannel(ctx, clientId, userKey)
	// if err != nil {
	// 	l.Error("could not register ws amqpChannel", zap.Error(err))
	// 	return err
	// }
	//
	// err := h.groupService.RegisterUserAmqpSubscriptions(ctx)
	//
	// c.Logger().Info("getting subscriptions")
	// getSubscriptions := chat.NewGetSubscriptions(userKey, 100, 0)
	// subs, err := h.chatStore.GetSubscriptionsForUser(ctx, getSubscriptions)
	// if err != nil {
	// 	l.Error("could not get subscriptions", zap.Error(err))
	// 	return err
	// }
	//
	// c.Logger().Info("creating subscriptions for user")
	// for _, item := range subs.Subscriptions.Items {
	// 	err := h.mq.BindUserExchangeToChannel(ctx, item.GetChannelKey(), userKey)
	// 	if err != nil {
	// 		l.Error("could not create rabbitmq binding for subscription",
	// 			zap.String("channel_id", item.ChannelID),
	// 			zap.String("user_id", item.UserID),
	// 			zap.Error(err),
	// 		)
	// 		return err
	// 	}
	// }
	//
	// membershipStatus := group.ApprovedMembershipStatus
	// getMemberships := group.NewGetMembershipsForUserRequest(userKey, &membershipStatus)
	// getMembershipsResponse := h.groupStore.GetMembershipsForUser(getMemberships)
	// if getMembershipsResponse.Error != nil {
	// 	l.Error("could not get group memberships for user", zap.Error(err))
	// 	return getMembershipsResponse.Error
	// }
	//
	// for _, membership := range getMembershipsResponse.Memberships.Items {
	// 	err = h.mq.RegisterUserMembershipBinding(membership.GetUserKey(), membership.GetGroupKey())
	// 	if err != nil {
	// 		l.Error("could not register user membership binding", zap.Error(err))
	// 		return err
	// 	}
	// }

}
