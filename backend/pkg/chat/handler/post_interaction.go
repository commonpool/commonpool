package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/auth"
	model2 "github.com/commonpool/backend/pkg/chat/handler/model"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
	"github.com/commonpool/backend/pkg/handler"
	usermodel "github.com/commonpool/backend/pkg/user/model"
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
	loggedInUserKey := usermodel.NewUserKey(loggedInUser.Subject)

	req := model2.SubmitInteractionRequest{}
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

	message, err := chatHandler.chatService.GetMessage(ctx, chatmodel.NewMessageKey(uid))
	if err != nil {
		return err
	}

	now := time.Now()

	var actions []model2.Action
	for _, action := range req.Payload.Actions {
		actions = append(actions, model2.Action{
			SubmitAction: model2.SubmitAction{
				ElementState: action.ElementState,
				BlockID:      action.BlockID,
				ActionID:     action.ActionID,
			},
			ActionTimestamp: now,
		})
	}

	webMessage := model2.MapMessage(message)

	interactionPayload := model2.InteractionCallback{
		Token: chatHandler.appConfig.CallbackToken,
		Payload: model2.InteractionCallbackPayload{
			Type:        model2.BlockActions,
			TriggerId:   "",
			ResponseURL: "",
			User: model2.InteractionPayloadUser{
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
