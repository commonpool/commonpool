package handler

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"net/http"
)

func (s *chatHandlerSuite) TestPostMessage() {
	s.NewContext(http.MethodPost, "/api/v1/chat/channel-1", `{"message":"hello"}`)
	s.LoggedInAs("user", "username", "user@email.com")
	s.ChatService.SendMessageFunc = func(ctx context.Context, message *chat.Message) error {
		return nil
	}
	s.ServeHTTP()
	s.AssertAccepted()
}

func (s *chatHandlerSuite) TestPostMessageShouldFailIfMessageNil() {
	s.NewContext(http.MethodPost, "/api/v1/chat/channel-1", `{"message":null}`)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		HasStatusCode(http.StatusBadRequest),
		HasCode("ErrValidation"),
		HasValidationError("SendMessageRequest.message", "message is required"))
}

func (s *chatHandlerSuite) TestPostMessageShouldFailIfMessageEmpty() {
	s.NewContext(http.MethodPost, "/api/v1/chat/channel-1", `{"message":""}`)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		HasStatusCode(http.StatusBadRequest),
		HasCode("ErrValidation"),
		HasValidationError("SendMessageRequest.message", "message is required"))
}

func (s *chatHandlerSuite) TestPostMessageShouldFailIfMessageBlank() {
	s.NewContext(http.MethodPost, "/api/v1/chat/channel-1", `{"message":"   "}`)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		HasStatusCode(http.StatusBadRequest),
		HasCode("ErrValidation"),
		HasValidationError("SendMessageRequest.message", "message is required"))
}
