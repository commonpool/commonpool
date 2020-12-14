package handler

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/logging"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func GetContext(c echo.Context) context.Context {
	ctx := c.Request().Context()

	if isAuth, ok := c.Get(auth.IsAuthenticatedKey).(bool); ok {
		if isAuth {
			ctx = context.WithValue(ctx, auth.IsAuthenticatedKey, true)
			ctx = context.WithValue(ctx, auth.SubjectUsernameKey, c.Get(auth.SubjectUsernameKey))
			ctx = context.WithValue(ctx, auth.SubjectEmailKey, c.Get(auth.SubjectEmailKey))
			ctx = context.WithValue(ctx, auth.SubjectKey, c.Get(auth.SubjectKey))
		}
	}

	return ctx
}

func GetCtx(ctx context.Context, handler string) (context.Context, *zap.Logger) {
	l := logging.WithContext(ctx).With(zap.String("handler", handler)).Named("handler." + handler)
	return ctx, l
}
