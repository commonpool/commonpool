package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

var requestMu = sync.Mutex{}

func NewRequest(ctx context.Context, session *models.UserSession, method, target string, req interface{}) (*http.Request, *httptest.ResponseRecorder) {
	requestMu.Lock()
	httpRequest := httptest.NewRequest(method, target, read(req))
	httpRequest = httpRequest.WithContext(ctx)
	httpRequest.Header.Set("Content-Type", "application/json")
	if session != nil {
		httpRequest.Header.Set("X-Debug-Username", session.Username)
		httpRequest.Header.Set("X-Debug-Email", session.Email)
		httpRequest.Header.Set("X-Debug-User-Id", session.Subject)
		httpRequest.Header.Set("X-Debug-Is-Authenticated", strconv.FormatBool(session.IsAuthenticated))
	}
	recorder := httptest.NewRecorder()
	requestMu.Unlock()
	return httpRequest, recorder
}

func ReadResponse(t *testing.T, recorder *httptest.ResponseRecorder, output interface{}) *http.Response {
	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, output)
	indent, err := json.MarshalIndent(output, "", "  ")
	if err == nil {
		t.Log("\n" + string(indent))
	}
	return resp
}

func read(intf interface{}) io.Reader {
	bts, err := json.Marshal(intf)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bts)
}

func setAuthenticatedUser(c echo.Context, username, subject, email string) {
	c.Set(authenticator.IsAuthenticatedKey, true)
	c.Set(authenticator.SubjectUsernameKey, username)
	c.Set(authenticator.SubjectEmailKey, email)
	c.Set(authenticator.SubjectKey, subject)
}
func setUnauthenticated(c echo.Context) {
	c.Set(authenticator.IsAuthenticatedKey, false)
}

func (s *IntegrationTestSuite) ListenOnUserExchange(t *testing.T, ctx context.Context, userKey keys.UserKey) *UserExchangeListener {

	randomStr := uuid.NewV4().String()
	amqpChan, err := s.server.AmqpClient.GetChannel()
	assert.NoError(s.T(), err)
	_, err = s.server.ChatService.CreateUserExchange(ctx, userKey)
	assert.NoError(s.T(), err)
	err = amqpChan.QueueDeclare(ctx, randomStr, false, true, false, false, nil)
	assert.NoError(s.T(), err)
	err = amqpChan.QueueBind(ctx, randomStr, "", userKey.GetExchangeName(), false, nil)
	assert.NoError(s.T(), err)
	delivery, err := amqpChan.Consume(ctx, randomStr, randomStr, false, false, false, false, nil)
	assert.NoError(s.T(), err)

	del := make(chan mq.Delivery)

	go func() {
		for d := range delivery {
			del <- d
		}
		close(del)
	}()

	return &UserExchangeListener{
		Channel:  amqpChan,
		Delivery: del,
	}

}

type UserExchangeListener struct {
	Channel  mq.Channel
	Delivery <-chan mq.Delivery
}

func (l *UserExchangeListener) Close() error {
	return l.Channel.Close()
}

func AssertStatusCreated(t *testing.T, httpResponse *http.Response) bool {
	return assert.NotNil(t, httpResponse) && assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
}

func AssertStatusNoContent(t *testing.T, httpResponse *http.Response) bool {
	return assert.NotNil(t, httpResponse) && assert.Equal(t, http.StatusNoContent, httpResponse.StatusCode)
}

func AssertOK(t *testing.T, httpResponse *http.Response) bool {
	return assert.NotNil(t, httpResponse) && assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
}

func AssertStatusUnauthorized(t *testing.T, httpResponse *http.Response) bool {
	return assert.NotNil(t, httpResponse) && assert.Equal(t, http.StatusUnauthorized, httpResponse.StatusCode)
}
func AssertStatusForbidden(t *testing.T, httpResponse *http.Response) bool {
	return assert.NotNil(t, httpResponse) && assert.Equal(t, http.StatusForbidden, httpResponse.StatusCode)
}

func AssertStatusBadRequest(t *testing.T, httpResponse *http.Response) bool {
	return assert.NotNil(t, httpResponse) && assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
}

func AssertStatusAccepted(t *testing.T, httpResponse *http.Response) bool {
	return assert.NotNil(t, httpResponse) && assert.Equal(t, http.StatusAccepted, httpResponse.StatusCode)
}
