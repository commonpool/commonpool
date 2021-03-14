package realtime

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/chat/service"
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
	chatService   service.Service
	authorization authenticator.Authenticator
	mq            mq.MqClient
}

func NewRealtimeHandler(mq mq.MqClient, chatService service.Service, authorization authenticator.Authenticator) *Handler {
	return &Handler{
		mq:            mq,
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

	queue := h.mq.NewQueue(mq.NewQueueConfig().WithName("").WithExclusive(true).WithAutoDelete(true))
	if err := queue.Connect(ctx); err != nil {
		return err
	}

	client := NewAnonymousClient(hub, ws, queue)

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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	userSession, err := oidc.GetLoggedInUser(ctx)
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

	channel := h.mq.NewChannel()
	if err := channel.Connect(ctx); err != nil {
		return err
	}

	defer channel.Close()

	userExchangeName, err := h.chatService.CreateUserExchange(ctx, userKey)
	if err != nil {
		l.Error("could not create user exchange", zap.Error(err))
		return err
	}

	queue := h.mq.NewQueue(mq.NewQueueConfig().WithExchange(userExchangeName).WithName(queueName).WithExclusive(false).WithAutoDelete(true))
	if err := queue.Connect(ctx); err != nil {
		return err
	}
	defer queue.Close()

	if err := channel.QueueBind(ctx, queueName, "", userExchangeName, false, nil); err != nil {
		l.Error("could not bind consumer queue to exchange", zap.Error(err))
		return err
	}

	client := NewClient(hub, ws, queue, &queueName, &consumerKey)
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
	// 			zap.String("user_id", item.UserKey),
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
