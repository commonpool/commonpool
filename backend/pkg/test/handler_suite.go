package test

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

type HandlerSuite struct {
	suite.Suite
	Echo          *echo.Echo
	Request       *http.Request
	Recorder      *httptest.ResponseRecorder
	Context       echo.Context
	Authenticator *mock.AuthenticatorMock
	LoggedInUser  models.UserReference
}

func (s *HandlerSuite) SetupTest() {
	s.Echo = echo.New()
	s.Echo.Validator = validation.DefaultValidator
	s.Echo.HTTPErrorHandler = handler.HttpErrorHandler
	s.Authenticator = &mock.AuthenticatorMock{
		AuthenticateFunc: func(redirectOnError bool) echo.MiddlewareFunc {
			return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
				return func(context echo.Context) error {
					return handlerFunc(context)
				}
			}
		},
		GetLoggedInUserFunc: func(ctx context.Context) (models.UserReference, error) {
			if s.LoggedInUser != nil {
				return s.LoggedInUser, nil
			}
			return nil, exceptions.ErrUnauthorized
		},
		GetRedirectResponseFunc: nil,
		LoginFunc:               nil,
		LogoutFunc:              nil,
	}
	s.Recorder = httptest.NewRecorder()
}

func (s *HandlerSuite) SetQueryParamString(ctx echo.Context, key string, value string) {
	ctx.QueryParams().Set(key, value)
}

func (s *HandlerSuite) SetQueryParamTimestamp(ctx echo.Context, key string, value time.Time) {
	s.SetQueryParamInt64(ctx, key, value.Unix())
}

func (s *HandlerSuite) SetQueryParamInt(ctx echo.Context, key string, value int) {
	ctx.QueryParams().Set(key, strconv.Itoa(value))
}

func (s *HandlerSuite) SetQueryParamInt64(ctx echo.Context, key string, value int64) {
	s.SetQueryParamInt(ctx, key, int(value))
}

func (s *HandlerSuite) NewContext(method string, target string, body interface{}) echo.Context {
	var reader io.Reader = nil
	if body != nil {
		if _, ok := body.(string); ok {
			reader = strings.NewReader(body.(string))
		} else {
			js, err := json.Marshal(body)
			if err != nil {
				s.T().Fatal(err)
			}
			reader = strings.NewReader(string(js))
		}
	}
	req := httptest.NewRequest(method, target, reader)
	if reader != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	s.Request = req
	c := s.Echo.NewContext(req, s.Recorder)
	s.Context = c
	return c
}

func (s *HandlerSuite) LoggedInAs(userKey string, username string, email string) {
	s.LoggedInUser = &models.UserSession{
		Username:        username,
		Subject:         userKey,
		Email:           email,
		IsAuthenticated: true,
	}
}

func (s *HandlerSuite) ReadResponse(response interface{}) {
	bodyBytes, err := ioutil.ReadAll(s.Recorder.Body)
	if !assert.Nil(s.T(), err) {
		s.T().Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &response)
	if !assert.Nil(s.T(), err) {
		s.T().Fatal(err)
	}
}

func (s *HandlerSuite) ServeHTTP() {
	s.Echo.ServeHTTP(s.Recorder, s.Request)
}

func (s *HandlerSuite) AssertResponseCode(expected int) bool {
	return assert.Equal(s.T(), expected, s.Recorder.Code)
}

func (s *HandlerSuite) AssertOK() bool {
	return s.AssertResponseCode(http.StatusOK)
}

func (s *HandlerSuite) AssertBadRequest() bool {
	return s.AssertResponseCode(http.StatusBadRequest)
}

func (s *HandlerSuite) AssertAccepted() bool {
	return s.AssertResponseCode(http.StatusAccepted)
}

func (s *HandlerSuite) AssertNoContent() bool {
	return s.AssertResponseCode(http.StatusNoContent)
}

func (s *HandlerSuite) AssertCreated() bool {
	return s.AssertResponseCode(http.StatusCreated)
}

type ErrorResponseAssertion = func(t *testing.T) func(err *exceptions.ErrorResponse) bool

func (s *HandlerSuite) AssertErrorResponse(do ...ErrorResponseAssertion) bool {
	wse := exceptions.ErrorResponse{}
	if err := json.Unmarshal(s.Recorder.Body.Bytes(), &wse); err != nil {
		s.T().Fatal(err)
		return false
	}
	for _, doFunc := range do {
		if !doFunc(s.T())(&wse) {
			return false
		}
	}
	return true
}

func HasStatusCode(statusCode int) ErrorResponseAssertion {
	return func(t *testing.T) func(err *exceptions.ErrorResponse) bool {
		return func(err *exceptions.ErrorResponse) bool {
			return assert.Equal(t, statusCode, err.StatusCode)
		}
	}
}

func HasCode(code string) ErrorResponseAssertion {
	return func(t *testing.T) func(err *exceptions.ErrorResponse) bool {
		return func(err *exceptions.ErrorResponse) bool {
			return assert.Equal(t, code, err.Code)
		}
	}
}

func HasMessage(message string) ErrorResponseAssertion {
	return func(t *testing.T) func(err *exceptions.ErrorResponse) bool {
		return func(err *exceptions.ErrorResponse) bool {
			return assert.Equal(t, message, err.Message)
		}
	}
}

func HasValidationError(key string, message string) ErrorResponseAssertion {
	return func(t *testing.T) func(err *exceptions.ErrorResponse) bool {
		return func(err *exceptions.ErrorResponse) bool {
			if !assert.NotNil(t, err, "error should not be nil") {
				return false
			}
			if !assert.NotNil(t, err.Validation, "validation errors should not be nil") {
				return false
			}
			msg, ok := err.Validation[key]
			if !assert.True(t, ok, "validation errors should contain the %s key", key) {
				return false
			}
			return assert.Equal(t, message, msg)
		}
	}
}
