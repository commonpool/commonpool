package module

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/auth/handler"
	"github.com/commonpool/backend/pkg/auth/queries"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Queries struct {
	GetUser     *queries.GetUser
	GetUsername *queries.GetUsername
}

type Module struct {
	Handler        *handler.UserHandler
	Authenticator  authenticator.Authenticator
	Queries        Queries
	UserRepository domain.UserRepository
}

func NewUserModule(appConfig *config.AppConfig, db *gorm.DB, eventStore eventstore.EventStore) *Module {
	userRepository := store.NewEventSourcedUserRepository(eventStore)
	authenticator := oidc.NewAuth(appConfig, "/api/v1", userRepository)
	searchUsers := queries.NewSearchUsers(db)
	getUser := queries.NewGetUser(db)
	userHandler := handler.NewUserHandler(authenticator, getUser, searchUsers)

	return &Module{
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
