package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/commonpool/backend/logging"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const state = "someState"

func newState(c echo.Context, state string) (string, error) {

	l := logging.WithContext(c.Request().Context())

	origin := c.Request().Header.Get("Origin")
	referrer := c.Request().Referer()

	l.Debug("encoding nonce", zap.String("origin", origin), zap.String("referrer", referrer))

	st := Nonce{
		DesiredUrl: referrer,
		State:      state,
	}
	bytes, err := json.Marshal(st)
	if err != nil {
		l.Error("could not encode state", zap.Error(err))
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)
	return b64, nil
}

func decodeState(ctx context.Context, state string) (*Nonce, error) {

	l := logging.WithContext(ctx)

	bytes, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		l.Error("could not decode state", zap.Error(err))
		return nil, err
	}
	nonce := &Nonce{}
	err = json.Unmarshal(bytes, nonce)
	if err != nil {
		l.Error("could not unmarshal state", zap.Error(err))
		return nil, err
	}
	return nonce, nil
}

type Nonce struct {
	DesiredUrl string `json:"des"`
	State      string `json:"state"`
}
