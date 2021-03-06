package handler

import (
	"context"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func GetContext(c echo.Context) context.Context {
	ctx := c.Request().Context()

	if isAuth, ok := c.Get(authenticator.IsAuthenticatedKey).(bool); ok {
		if isAuth {
			ctx = context.WithValue(ctx, authenticator.IsAuthenticatedKey, true)
			ctx = context.WithValue(ctx, authenticator.SubjectUsernameKey, c.Get(authenticator.SubjectUsernameKey))
			ctx = context.WithValue(ctx, authenticator.SubjectEmailKey, c.Get(authenticator.SubjectEmailKey))
			ctx = context.WithValue(ctx, authenticator.SubjectKey, c.Get(authenticator.SubjectKey))
		}
	}

	return ctx
}

func GetCtx(ctx context.Context, handler string) (context.Context, *zap.Logger) {
	l := logging.WithContext(ctx).With(zap.String("handler", handler)).Named("handler." + handler)
	return ctx, l
}
