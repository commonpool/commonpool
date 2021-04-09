package main

import (
	"cp/pkg/acknowledgements"
	"cp/pkg/api"
	"cp/pkg/credits"
	"cp/pkg/groups"
	"cp/pkg/handler"
	"cp/pkg/images"
	"cp/pkg/memberships"
	"cp/pkg/messages"
	"cp/pkg/notifications"
	"cp/pkg/posts"
	"cp/pkg/users"
	"cp/pkg/utils"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"io"
	"os"
	"time"
)

type TemplateRenderer struct {
	templates         *template.Template
	cookieStore       *sessions.CookieStore
	membershipStore   memberships.Store
	alertManager      *utils.AlertManager
	notificationStore notifications.Store
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["path"] = c.Path()

		viewContext["Alerts"] = []utils.Alert{}
		alerts, err := t.alertManager.GetAlerts(c.Request())
		if err == nil {
			viewContext["Alerts"] = alerts
			_ = t.alertManager.ClearAlerts(c.Request(), c.Response().Writer)
		}

		viewName, err := utils.GetRouteName(c)
		if err != nil {
			return err
		}
		viewContext["route"] = viewName

		profile, err := handler.GetProfile(t.cookieStore, c)
		if err != nil {
			c.Logger().Error(fmt.Errorf("failed to get profile: %w", err))
			return echo.ErrInternalServerError
		}

		if profile == nil {
			viewContext["IsAuthenticated"] = false
		} else {

			var adminMemberships []*api.Membership
			permission := api.Admin
			_ = t.membershipStore.Find(&adminMemberships, &memberships.GetMembershipsOptions{
				HasPermission: &permission,
				UserID:        &profile.ID,
				Preload:       []string{"Group"},
			})

			var administeredGroups []*api.Group
			for _, membership := range adminMemberships {
				administeredGroups = append(administeredGroups, membership.Group)
			}

			viewContext["IsAuthenticated"] = true
			viewContext["Session"] = map[string]interface{}{
				"Email":              profile.Email,
				"Username":           profile.Username,
				"UserID":             profile.ID,
				"AdministeredGroups": administeredGroups,
			}

			count, err := t.notificationStore.GetUnreadCount(profile.ID)
			if err != nil {
				return err
			}

			viewContext["unreadNotificationCount"] = count

		}

		return t.templates.Funcs(map[string]interface{}{
			"session": func() interface{} {
				return viewContext["Session"]
			},
			"viewName": func() string {
				return viewName
			},
			"isView": func(n string) bool {
				return n == viewName
			},
			"html": func(s string) template.HTML {
				return template.HTML(s)
			},
			"getMembership": func(groupID string, userID string) (*api.Membership, error) {
				if userID == "" {
					return nil, nil
				}
				if groupID == "" {
					return nil, nil
				}
				m, err := t.membershipStore.Get(groupID, userID)
				if errors.Is(err, echo.ErrNotFound) {
					return nil, nil
				}
				if err != nil {
					return nil, err
				}
				return m, nil
			},
			"loggedInUserID": func() string {
				if userID, ok := c.Get("loggedInUserID").(string); ok {
					return userID
				}
				return ""
			},
			"loggedInUsername": func() string {
				if username, ok := c.Get("loggedInUsername").(string); ok {
					return username
				}
				return ""
			},
			"groupID": func() string {
				if username, ok := c.Get("groupID").(string); ok {
					return username
				}
				return ""
			},
			"json": func(v interface{}) (string, error) {
				bytes, err := json.Marshal(v)
				if err != nil {
					return "", err
				}
				return string(bytes), nil
			},
			handler.GroupKey: func() *api.Group {
				g := c.Get(handler.GroupKey)
				if g == nil {
					return nil
				}
				return g.(*api.Group)
			},
			handler.GroupIDKey: func() string {
				g := c.Get(handler.GroupIDKey)
				if g == nil {
					return ""
				}
				return g.(string)
			},
			handler.UserKey: func() *api.User {
				u := c.Get(handler.UserKey)
				if u == nil {
					return nil
				}
				return u.(*api.User)
			},
			handler.UserIDKey: func() string {
				u := c.Get(handler.UserIDKey)
				if u == nil {
					return ""
				}
				return u.(string)
			},
			handler.MembershipKey: func() *api.Membership {
				m := c.Get(handler.MembershipKey)
				if m == nil {
					return nil
				}
				return m.(*api.Membership)
			},
			handler.PostIDKey: func() string {
				postID := c.Get(handler.PostIDKey)
				if postID == nil {
					return ""
				}
				return postID.(string)
			},
			handler.PostKey: func() *api.Post {
				p := c.Get(handler.PostKey)
				if p == nil {
					return nil
				}
				return p.(*api.Post)
			},
			handler.AuthenticatedUserKey: func() *api.User {
				u := c.Get(handler.AuthenticatedUserKey)
				if u == nil {
					return nil
				}
				return u.(*api.User)
			},
			handler.AuthenticatedUserMembershipKey: func() *api.Membership {
				m := c.Get(handler.AuthenticatedUserMembershipKey)
				if m == nil {
					return nil
				}
				return m.(*api.Membership)
			},
			handler.ProfileKey: func() *api.Profile {
				p := c.Get(handler.ProfileKey)
				if p == nil {
					return nil
				}
				return p.(*api.Profile)
			},
		}).ExecuteTemplate(w, name, data)

	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	gob.Register([]utils.Alert{})

	var database *gorm.DB
	var err error

	dbProvider := os.Getenv("DB_PROVIDER")
	if dbProvider == "" || dbProvider == "sqlite" {
		database, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		database.DisableForeignKeyConstraintWhenMigrating = true
	} else if dbProvider == "postgres" {

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)

		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		db, err := database.DB()
		if err != nil {
			panic(err)
		}
		db.SetConnMaxLifetime(time.Minute * 5)
		db.SetConnMaxIdleTime(time.Minute * 5)

	}

	if err := database.AutoMigrate(
		&api.Group{},
		&api.Membership{},
		&api.User{},
		&api.Post{},
		&api.Message{},
		&api.Acknowledgement{},
		&api.Credits{},
		&api.Notification{},
		&api.Image{},
	); err != nil {
		panic(err)
	}

	groupStore := groups.NewGroupStore(database)
	membershipStore := memberships.NewMembershipStore(database)
	userStore := users.NewUserStore(database)
	postStore := posts.NewPostStore(database)
	messageStore := messages.NewMessageStore(database)
	acknowledgementStore := acknowledgements.NewAcknowledgementStore(database)
	creditsStore := credits.NewCreditStore(database)
	cookieStore := sessions.NewCookieStore([]byte("secret"))
	alertManager := utils.NewAlertManager(cookieStore)
	notificationStore := notifications.NewNotificationStore(database)
	imageStore := images.NewImageStore(database)

	_, _ = template.New("").Funcs(map[string]interface{}{
		"session": func() interface{} {
			return nil
		},
	}).Parse("")

	funcMap := map[string]interface{}{
		"session": func() interface{} {
			return nil
		},
		"viewName": func() string {
			return ""
		},
		"isView": func(vn string) bool {
			return false
		},
		"html": func(html string) template.HTML {
			return template.HTML("")
		},
		"getMembership": func(groupID string, userID string) (*api.Membership, error) {
			return nil, nil
		},
		"loggedInUserID": func() string {
			return ""
		},
		"loggedInUsername": func() string {
			return ""
		},
		"json": func(v interface{}) (string, error) {
			return "", nil
		},
		handler.GroupKey: func() *api.Group {
			return nil
		},
		handler.GroupIDKey: func() string {
			return ""
		},
		handler.UserKey: func() *api.User {
			return nil
		},
		handler.UserIDKey: func() string {
			return ""
		},
		handler.MembershipKey: func() *api.Membership {
			return nil
		},
		handler.PostIDKey: func() string {
			return ""
		},
		handler.PostKey: func() *api.Post {
			return nil
		},
		handler.AuthenticatedUserKey: func() *api.User {
			return nil
		},
		handler.AuthenticatedUserMembershipKey: func() *api.Membership {
			return nil
		},
		handler.ProfileKey: func() *api.Profile {
			return nil
		},
	}

	viewsDir := os.Getenv("VIEWS_DIR")
	if viewsDir == "" {
		viewsDir = "public/views"
	}

	renderer := &TemplateRenderer{
		templates: template.Must(
			template.New("main").Funcs(funcMap).ParseGlob(fmt.Sprintf("%s/*.gohtml", viewsDir)),
		),
		cookieStore:       cookieStore,
		membershipStore:   membershipStore,
		alertManager:      alertManager,
		notificationStore: notificationStore,
	}

	h := handler.NewHandler(
		cookieStore,
		groupStore,
		membershipStore,
		userStore,
		postStore,
		creditsStore,
		acknowledgementStore,
		messageStore,
		notificationStore,
		imageStore,
		alertManager,
		database,
	)

	e := echo.New()
	e.Renderer = renderer
	e.Debug = true
	e.Use(session.Middleware(cookieStore))

	h.Register(e)

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = ":8000"
	}
	e.Logger.Fatal(e.Start(listenAddress))
}
