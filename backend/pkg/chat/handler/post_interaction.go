package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

// SubmitInteraction
// @Summary Sends a message to a topic
// @Description This endpoint is for user interactions through the chat box
// @ID submitInteraction
// @Param message body web.SubmitInteractionRequest true "Message to send"
// @Tags chat
// @Accept json
// @Success 200
// @Failure 400 {object} utils.Error
// @Router /chat/interaction [post]
func (chatHandler *ChatHandler) SubmitInteraction(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "SubmitInteraction")

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := model.NewUserKey(loggedInUser.Subject)

	req := web.SubmitInteractionRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	uid, err := uuid.FromString(req.Payload.MessageID)
	if err != nil {
		return err
	}

	message, err := chatHandler.chatService.GetMessage(ctx, model.NewMessageKey(uid))
	if err != nil {
		return err
	}

	now := time.Now()

	var actions []web.Action
	for _, action := range req.Payload.Actions {
		actions = append(actions, web.Action{
			SubmitAction: web.SubmitAction{
				ElementState: action.ElementState,
				BlockID:      action.BlockID,
				ActionID:     action.ActionID,
			},
			ActionTimestamp: now,
		})
	}

	webMessage := web.MapMessage(message)

	interactionPayload := web.InteractionCallback{
		Token: chatHandler.appConfig.CallbackToken,
		Payload: web.InteractionCallbackPayload{
			Type:        web.BlockActions,
			TriggerId:   "",
			ResponseURL: "",
			User: web.InteractionPayloadUser{
				ID:       loggedInUserKey.String(),
				Username: loggedInUser.Username,
			},
			Message: webMessage,
			Actions: actions,
			State:   req.Payload.State,
		},
	}

	requestBody, err := json.Marshal(interactionPayload)
	if err != nil {
		return err
	}

	httpRequest, err := http.NewRequest("POST", "http://localhost:8585/api/v1/chatback", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+c.Get("token").(string))

	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code")
	}

	return c.String(http.StatusOK, "OK")

}
