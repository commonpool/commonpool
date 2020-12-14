package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/router"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var requestMu = sync.Mutex{}

func NewRequest(ctx context.Context, session *auth.UserSession, method, target string, req interface{}) (echo.Context, *httptest.ResponseRecorder) {
	requestMu.Lock()
	e := router.NewRouter()
	httpRequest := httptest.NewRequest(method, target, read(req))
	httpRequest = httpRequest.WithContext(ctx)
	httpRequest.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c := e.NewContext(httpRequest, recorder)
	if session != nil {
		setAuthenticatedUser(c, session.Username, session.Subject, session.Email)
	} else {
		setUnauthenticated(c)
	}
	requestMu.Unlock()
	return c, recorder
}

func ReadResponse(t *testing.T, recorder *httptest.ResponseRecorder, output interface{}) *http.Response {
	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, output)
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
	c.Set(auth.IsAuthenticatedKey, true)
	c.Set(auth.SubjectUsernameKey, username)
	c.Set(auth.SubjectEmailKey, email)
	c.Set(auth.SubjectKey, subject)
}
func setUnauthenticated(c echo.Context) {
	c.Set(auth.IsAuthenticatedKey, false)
}

func ListenOnUserExchange(t *testing.T, ctx context.Context, userKey model.UserKey) *UserExchangeListener {

	randomStr := uuid.NewV4().String()
	amqpChan, err := AmqpClient.GetChannel()
	assert.NoError(t, err)
	_, err = ChatService.CreateUserExchange(ctx, userKey)
	assert.NoError(t, err)
	err = amqpChan.QueueDeclare(ctx, randomStr, false, true, false, false, nil)
	assert.NoError(t, err)
	err = amqpChan.QueueBind(ctx, randomStr, "", userKey.GetExchangeName(), false, nil)
	assert.NoError(t, err)
	delivery, err := amqpChan.Consume(ctx, randomStr, randomStr, false, false, false, false, nil)
	assert.NoError(t, err)

	del := make(chan amqp.Delivery)

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
	Channel  amqp.Channel
	Delivery <-chan amqp.Delivery
}

func (l *UserExchangeListener) Close() error {
	return l.Channel.Close()
}

func AssertStatusCreated(t *testing.T, httpResponse *http.Response) {
	assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
}

func AssertStatusNoContent(t *testing.T, httpResponse *http.Response) {
	assert.Equal(t, http.StatusNoContent, httpResponse.StatusCode)
}

func AssertOK(t *testing.T, httpResponse *http.Response) {
	assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
}
