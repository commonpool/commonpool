package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bsm/redislock"
	_ "github.com/commonpool/backend/docs"
	commands2 "github.com/commonpool/backend/pkg/auth/commands"
	authdomain "github.com/commonpool/backend/pkg/auth/domain"
	userlisteners "github.com/commonpool/backend/pkg/auth/listeners"
	"github.com/commonpool/backend/pkg/auth/module"
	chathandler "github.com/commonpool/backend/pkg/chat/handler"
	chatservice "github.com/commonpool/backend/pkg/chat/service"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/clusterlock"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/config"
	db2 "github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	postgres2 "github.com/commonpool/backend/pkg/eventstore/postgres"
	"github.com/commonpool/backend/pkg/eventstore/publish"
	"github.com/commonpool/backend/pkg/graph"
	groupdomain "github.com/commonpool/backend/pkg/group/domain"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	grouplisteners "github.com/commonpool/backend/pkg/group/listeners"
	groupqueries "github.com/commonpool/backend/pkg/group/queries"
	groupservice "github.com/commonpool/backend/pkg/group/service"
	groupstore "github.com/commonpool/backend/pkg/group/store"
	handler2 "github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/mq"
	nukehandler "github.com/commonpool/backend/pkg/nuke/handler"
	"github.com/commonpool/backend/pkg/realtime"
	resourcedomain "github.com/commonpool/backend/pkg/resource/domain"
	resourcehandler "github.com/commonpool/backend/pkg/resource/handler"
	resourcelisteners "github.com/commonpool/backend/pkg/resource/listeners"
	resourcequeries "github.com/commonpool/backend/pkg/resource/queries"
	resourcestore "github.com/commonpool/backend/pkg/resource/store"
	"github.com/commonpool/backend/pkg/trading"
	tradingcmdhandlers "github.com/commonpool/backend/pkg/trading/commandhandlers"
	tradingdomain "github.com/commonpool/backend/pkg/trading/domain"
	tradinghandler "github.com/commonpool/backend/pkg/trading/handler"
	tradinglisteners "github.com/commonpool/backend/pkg/trading/listeners"
	"github.com/commonpool/backend/pkg/trading/queries"
	tradingservice "github.com/commonpool/backend/pkg/trading/service"
	tradingstore "github.com/commonpool/backend/pkg/trading/store"
	"github.com/commonpool/backend/pkg/transaction"
	transactionservice "github.com/commonpool/backend/pkg/transaction/service"
	transactionstore "github.com/commonpool/backend/pkg/transaction/store"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"time"
)

type Group struct {
	Handler    *grouphandler.Handler
	Repository groupdomain.GroupRepository
	Service    *groupservice.GroupService
	Store      *groupstore.GroupStore
}

type Server struct {
	AppConfig                   *config.AppConfig
	AmqpClient                  mq.Client
	GraphDriver                 graph.Driver
	Db                          *gorm.DB
	TransactionStore            transaction.Store
	TransactionService          transaction.Service
	ChatStore                   chatstore.Store
	ChatService                 chatservice.Service
	TradingService              trading.Service
	ChatHandler                 *chathandler.Handler
	ResourceHandler             *resourcehandler.ResourceHandler
	RealTimeHandler             *realtime.Handler
	Router                      *echo.Echo
	NukeHandler                 *nukehandler.Handler
	TradingHandler              *tradinghandler.TradingHandler
	RedisClient                 *redis.Client
	ClusterLocker               clusterlock.Locker
	CommandMapper               *commands.CommandMapper
	CommandBus                  commands.CommandBus
	EventMapper                 *eventsource.EventMapper
	Group                       Group
	User                        *module.Module
	ErrGroup                    *errgroup.Group
	Ctx                         context.Context
	GetOffer                    *queries.GetOffer
	GetOffers                   *queries.GetOffers
	GetOfferItem                *queries.GetOfferItem
	GetGroup                    *groupqueries.GetGroup
	GetGroupByKeys              *groupqueries.GetGroupByKeys
	GetMembership               *groupqueries.GetMembershipReadModel
	GetOfferKeyForOfferItem     *queries.GetOfferKeyForOfferItemKey
	GetResource                 *resourcequeries.GetResource
	GetResourceWithSharings     *resourcequeries.GetResourceWithSharings
	GetResourcesByKeys          *resourcequeries.GetResourcesByKeys
	SearchResources             *resourcequeries.SearchResources
	SearchResourcesWithSharings *resourcequeries.SearchResourcesWithSharings
	AcceptOffer                 *tradingcmdhandlers.AcceptOfferHandler
	DeclineOffer                *tradingcmdhandlers.DeclineOfferHandler
	ConfirmResourceGiven        *tradingcmdhandlers.ConfirmResourceGivenHandler
	ConfirmResourceBorrowed     *tradingcmdhandlers.ConfirmResourceBorrowedHandler
	ConfirmServiceGiven         *tradingcmdhandlers.ConfirmServiceGivenHandler
	ConfirmResourceReturned     *tradingcmdhandlers.ConfirmResourceReturnedHandler
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func NewServer() (*Server, error) {

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	var redisTlsConfig *tls.Config = nil
	if appConfig.RedisTlsEnabled {
		redisTlsConfig = &tls.Config{
			InsecureSkipVerify: appConfig.RedisTlsSkipVerify,
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:      appConfig.RedisHost + ":" + appConfig.RedisPort,
		Password:  appConfig.RedisPassword,
		DB:        0,
		TLSConfig: redisTlsConfig,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	clusterLocker := clusterlock.NewRedis(redislock.New(redisClient))

	amqpCli, err := mq.NewRabbitMqClient(ctx, appConfig.AmqpUrl)
	if err != nil {
		return nil, err
	}

	err = graph.InitGraphDatabase(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	driver, err := graph.NewNeo4jDriver(appConfig, appConfig.Neo4jDatabase)
	if err != nil {
		return nil, err
	}

	db := getDb(appConfig)
	db2.AutoMigrate(db)

	// events
	eventMapper := eventsource.NewEventMapper()
	if err := authdomain.RegisterEvents(eventMapper); err != nil {
		panic(err)
	}
	if err := groupdomain.RegisterEvents(eventMapper); err != nil {
		panic(err)
	}
	if err := resourcedomain.RegisterEvents(eventMapper); err != nil {
		panic(err)
	}
	if err := tradingdomain.RegisterEvents(eventMapper); err != nil {
		panic(err)
	}

	eventPublisher := eventbus.NewAmqpPublisher(amqpCli)
	eventStore := publish.NewPublishEventStore(postgres2.NewPostgresEventStore(db, eventMapper), eventPublisher)
	// userRepo := store.NewEventSourcedUserRepository(eventStore)
	offerRepository := tradingstore.NewEventSourcedOfferRepository(eventStore)

	// queries
	getOfferKeyForOfferItemKeyQry := queries.NewGetOfferKeyForOfferItemKey(db)
	getGroupByKeys := groupqueries.NewGetGroupByKeys(db)
	getGroup := groupqueries.NewGetGroupReadModel(db)
	getMembership := groupqueries.NewGetMembership(db)
	getUserMemberships := groupqueries.NewGetUserMemberships(db)
	getGroupMemberships := groupqueries.NewGetGroupMemberships(db)
	getUsersForGroupInvite := groupqueries.NewGetUsersForGroupInvite(db)
	getResource := resourcequeries.NewGetResource(db)
	getResourcesByKeys := resourcequeries.NewGetResourcesByKeys(db)
	searchResources := resourcequeries.NewSearchResources(db)
	getResourceSharings := resourcequeries.NewGetResourceSharings(db)
	getResourcesSharings := resourcequeries.NewGetResourcesSharings(db)
	getResourceWithSharings := resourcequeries.NewGetResourceWithSharings(getResource, getResourceSharings)
	searchResourcesWithSharings := resourcequeries.NewSearchResourcesWithSharings(db, searchResources, getResourcesSharings)
	getOffer := queries.NewGetOffer(db)
	getOffers := queries.NewGetOffers(db)
	getOfferItem := queries.NewGetOfferItem(db)
	getOfferKeyForOfferItem := queries.NewGetOfferKeyForOfferItemKey(db)
	getOfferPermissions := queries.NewGetOfferPermissions(db)

	// commands
	commandMapper := commands.NewCommandMapper()
	// commandBus := rabbit.NewRabbitCommandBus(amqpCli, commandMapper)
	commands2.RegisterCommands(commandMapper)

	acceptOffer := tradingcmdhandlers.NewAcceptOfferHandler(offerRepository, getOfferPermissions)
	declineOffer := tradingcmdhandlers.NewDeclineOfferHandler(offerRepository)
	submitOffer := tradingcmdhandlers.NewSubmitOfferHandler(offerRepository, getOfferPermissions, getResourcesByKeys)
	confirmResourceBorrowed := tradingcmdhandlers.NewConfirmResourceBorrowedHandler(offerRepository)
	confirmResourceGiven := tradingcmdhandlers.NewConfirmResourceGivenHandler(offerRepository)
	confirmResourceReturned := tradingcmdhandlers.NewConfirmResourceReturnedHandler(offerRepository)
	confirmServiceGiven := tradingcmdhandlers.NewConfirmServiceGivenHandler(offerRepository)

	r := NewRouter()
	r.HTTPErrorHandler = handler2.HttpErrorHandler
	r.GET("/api/swagger/*", echoSwagger.WrapHandler)
	v1 := r.Group("/api/v1", middleware.Recover())

	userModule := module.NewUserModule(appConfig, db, driver, eventStore)
	userModule.Register(v1)

	transactionStore := transactionstore.NewTransactionStore(db)
	transactionService := transactionservice.NewTransactionService(transactionStore)

	resourceRepository := resourcestore.NewEventSourcedResourceRepository(eventStore)

	chatStore := chatstore.NewChatStore(db)
	chatService := chatservice.NewChatService(userModule.Store, amqpCli, chatStore)

	groupStore := groupstore.NewGroupStore(driver)
	groupRepo := groupstore.NewEventSourcedGroupRepository(eventStore)
	groupService := groupservice.NewGroupService(groupStore, amqpCli, chatService, userModule.Store, groupRepo, getGroup, getGroupByKeys, getGroupMemberships)

	tradingService := tradingservice.NewTradingService(groupService)

	chatHandler := chathandler.NewHandler(
		chatService,
		tradingService,
		appConfig,
		userModule.Authenticator,
		getOfferKeyForOfferItem,
		confirmResourceBorrowed,
		confirmResourceReturned,
		confirmServiceGiven,
		confirmResourceGiven,
		declineOffer,
		acceptOffer)
	chatHandler.Register(v1)

	groupHandler := grouphandler.NewHandler(
		groupService,
		userModule.Authenticator,
		getGroup,
		getMembership,
		getGroupMemberships,
		getUserMemberships,
		getUsersForGroupInvite)

	groupHandler.Register(v1)

	resourceHandler := resourcehandler.NewHandler(
		groupService,
		userModule.Service,
		userModule.Authenticator,
		resourceRepository,
		getUserMemberships,
		getResource,
		getResourceSharings,
		getResourcesSharings,
		searchResources,
		getResourceWithSharings,
		searchResourcesWithSharings)
	resourceHandler.Register(v1)

	realtimeHandler := realtime.NewRealtimeHandler(amqpCli, chatService, userModule.Authenticator)
	realtimeHandler.Register(v1)

	tradingHandler := tradinghandler.NewTradingHandler(
		tradingService,
		groupService,
		userModule.Service,
		userModule.Authenticator,
		getOffer,
		getOffers,
		getOfferKeyForOfferItem,
		confirmResourceBorrowed,
		confirmResourceReturned,
		confirmServiceGiven,
		confirmResourceGiven,
		declineOffer,
		acceptOffer,
		submitOffer,
	)

	tradingHandler.Register(v1)

	nukeHandler := nukehandler.NewHandler(db, amqpCli, driver)
	nukeHandler.Register(v1)

	var catchUpListenerFactory eventbus.CatchUpListenerFactory = func(key string, lockTTL time.Duration) *eventbus.CatchUpListener {
		return eventbus.NewCatchUpListener(
			eventStore,
			func() time.Time { return time.Time{} },
			amqpCli,
			eventbus.NewRedisDeduplicator(100, redisClient, key),
			clusterLocker,
			lockTTL,
			&clusterlock.Options{
				RetryStrategy: clusterlock.EverySecond,
			},
			eventMapper,
		)
	}

	handler := tradinglisteners.NewTransactionHistoryHandler(db, catchUpListenerFactory)
	g.Go(func() error {
		if err := handler.Start(ctx); err != nil {
			return err
		}
		return nil
	})

	offerRm := tradinglisteners.NewOfferReadModelHandler(db, catchUpListenerFactory)
	g.Go(func() error {
		if err := offerRm.Start(ctx); err != nil {
			return err
		}
		return nil
	})

	groupRm := grouplisteners.NewGroupReadModelListener(catchUpListenerFactory, db)
	g.Go(func() error {
		if err := groupRm.Start(ctx); err != nil {
			return err
		}
		return nil
	})

	userRm := userlisteners.NewUserReadModelListener(db, catchUpListenerFactory)
	g.Go(func() error {
		if err := userRm.Start(ctx); err != nil {
			return err
		}
		return nil
	})

	resourceRm := resourcelisteners.NewResourceReadModelHandler(db, catchUpListenerFactory)
	g.Go(func() error {
		if err := resourceRm.Start(ctx); err != nil {
			return err
		}
		return nil
	})

	return &Server{
		AppConfig:          appConfig,
		AmqpClient:         amqpCli,
		GraphDriver:        driver,
		Db:                 db,
		TransactionStore:   transactionStore,
		TransactionService: transactionService,
		ChatStore:          chatStore,
		ChatService:        chatService,
		TradingService:     tradingService,
		ChatHandler:        chatHandler,
		ResourceHandler:    resourceHandler,
		RealTimeHandler:    realtimeHandler,
		Router:             r,
		NukeHandler:        nukeHandler,
		TradingHandler:     tradingHandler,
		RedisClient:        redisClient,
		ClusterLocker:      clusterLocker,
		CommandMapper:      commandMapper,
		CommandBus:         nil,
		EventMapper:        eventMapper,
		Group: Group{
			Handler:    groupHandler,
			Service:    groupService,
			Repository: groupRepo,
			Store:      groupStore,
		},
		User:                        userModule,
		ErrGroup:                    g,
		Ctx:                         ctx,
		GetOffer:                    getOffer,
		GetOffers:                   getOffers,
		GetOfferItem:                getOfferItem,
		GetGroup:                    getGroup,
		GetGroupByKeys:              getGroupByKeys,
		GetMembership:               getMembership,
		GetOfferKeyForOfferItem:     getOfferKeyForOfferItemKeyQry,
		GetResource:                 getResource,
		GetResourceWithSharings:     getResourceWithSharings,
		SearchResources:             searchResources,
		SearchResourcesWithSharings: searchResourcesWithSharings,
		AcceptOffer:                 acceptOffer,
		DeclineOffer:                declineOffer,
		ConfirmResourceGiven:        confirmResourceGiven,
		ConfirmResourceBorrowed:     confirmResourceBorrowed,
		ConfirmServiceGiven:         confirmServiceGiven,
		ConfirmResourceReturned:     confirmResourceReturned,
	}, nil

}

func (s *Server) Start() error {
	if err := s.Router.Start("0.0.0.0:8585"); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.AmqpClient.Shutdown(); err != nil {
		return err
	}

	if err := s.Router.Shutdown(ctx); err != nil {
		return err
	}

	return nil

}

func getDb(appConfig *config.AppConfig) *gorm.DB {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", appConfig.DbHost, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName, appConfig.DbPort)
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
