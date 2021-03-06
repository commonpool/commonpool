package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
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
func (h *Handler) SubmitInteraction(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "SubmitInteraction")

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := keys.NewUserKey(loggedInUser.Subject)

	req := SubmitInteractionRequest{}
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

	message, err := h.chatService.GetMessage(ctx, keys.NewMessageKey(uid))
	if err != nil {
		return err
	}

	now := time.Now()

	var actions []Action
	for _, action := range req.Payload.Actions {
		actions = append(actions, Action{
			SubmitAction: SubmitAction{
				ElementState: action.ElementState,
				BlockID:      action.BlockID,
				ActionID:     action.ActionID,
			},
			ActionTimestamp: now,
		})
	}

	webMessage := MapMessage(message)

	interactionPayload := InteractionCallback{
		Token: h.appConfig.CallbackToken,
		Payload: InteractionCallbackPayload{
			Type:        BlockActions,
			TriggerId:   "",
			ResponseURL: "",
			User: InteractionPayloadUser{
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

type SubmitInteractionPayload struct {
	MessageID string                             `json:"messageId,omitempty"`
	State     map[string]map[string]ElementState `json:"state,omitempty"`
	Actions   []SubmitAction                     `json:"actions,omitempty"`
}

type SubmitInteractionRequest struct {
	Payload SubmitInteractionPayload `json:"payload"`
}

type SubmitAction struct {
	ElementState
	BlockID  string `json:"blockId,omitempty"`
	ActionID string `json:"actionId,omitempty"`
}

type ElementState struct {
	Type            chat.ElementType    `json:"type,omitempty"`
	SelectedDate    *string             `json:"selectedDate,omitempty"`
	SelectedTime    *string             `json:"selectedTime,omitempty"`
	Value           *string             `json:"value,omitempty"`
	SelectedOption  *chat.OptionObject  `json:"selectedOption,omitempty"`
	SelectedOptions []chat.OptionObject `json:"selectedOptions,omitempty"`
}

type Action struct {
	SubmitAction
	ActionTimestamp time.Time `json:"actionTimestamp,omitempty"`
}
