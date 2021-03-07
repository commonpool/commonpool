package module

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/auth/handler"
	"github.com/commonpool/backend/pkg/auth/queries"
	"github.com/commonpool/backend/pkg/auth/service"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/graph"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Queries struct {
	GetUser     *queries.GetUser
	GetUsername *queries.GetUsername
}

type Module struct {
	Store          store.Store
	Service        service.Service
	Handler        *handler.UserHandler
	Authenticator  authenticator.Authenticator
	Queries        Queries
	UserRepository domain.UserRepository
}

func NewUserModule(appConfig *config.AppConfig, db *gorm.DB, graphDriver graph.Driver, eventStore eventstore.EventStore) *Module {
	userRepository := store.NewEventSourcedUserRepository(eventStore)
	userStore := store.NewUserStore(graphDriver)
	userService := service.NewUserService(userStore)
	authenticator := oidc.NewAuth(appConfig, "/api/v1", userStore, userRepository)
	userHandler := handler.NewUserHandler(appConfig, userService, authenticator)

	return &Module{
		Store:          userStore,
		Service:        userService,
		Handler:        userHandler,
		Authenticator:  authenticator,
		UserRepository: userRepository,
		Queries: Queries{
			GetUser:     queries.NewGetUser(db),
			GetUsername: queries.NewGetUsername(db),
		},
	}
}

func (m *Module) Register(router *echo.Group) {
	m.Handler.Register(router)

}