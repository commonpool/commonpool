package handler

import (
	"cp/pkg/api"
	"cp/pkg/memberships"
	"cp/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
	"time"
)

type SendType string

const (
	Acknowledgement SendType = "acknowledgement"
	Credits         SendType = "credits"
	Other           SendType = "other"
)

func (s SendType) IsEmpty() bool {
	return string(s) == ""
}

func (s SendType) IsCredits() bool {
	return s == Credits
}
func (s SendType) IsAcknowledgement() bool {
	return s == Acknowledgement
}
func (s SendType) IsOther() bool {
	return s == Other
}

type SendOption struct {
	Type                SendType                `form:"type"`
	Target              string                  `form:"target"`
	Source              string                  `form:"source"`
	Amount              string                  `form:"amount"`
	AcknowledgementType api.AcknowledgementType `form:"acknowledgementType"`
	Notes               string                  `form:"notes"`
}

type SendTarget struct {
	DisplayName string
	Value       string
}

func (h *Handler) getTarget(targetStr string) (*api.Target, error) {

	if strings.Contains(targetStr, "user:") {

		userID := targetStr[5:]
		user, err := h.userStore.Get(userID)
		if err != nil {
			return nil, err
		}

		return &api.Target{
			UserID: &userID,
			User:   user,
			Type:   api.UserTarget,
		}, nil

	} else if strings.Contains(targetStr, "group:") {

		groupID := targetStr[6:]
		group, err := h.groupStore.Get(groupID)
		if err != nil {
			return nil, err
		}

		return &api.Target{
			GroupID: &groupID,
			Group:   group,
			Type:    api.GroupTarget,
		}, nil

	} else {
		return nil, echo.ErrBadRequest
	}
}

func (h *Handler) handleGroupSend(c echo.Context) error {

	authenticatedUser, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	group, err := h.getGroup(c)
	if err != nil {
		return err
	}

	var ms []*api.Membership
	if err := h.membershipStore.Find(&ms, &memberships.GetMembershipsOptions{
		GroupID: &group.ID,
		Preload: []string{"User"},
	}); err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet {

		var sources []*SendTarget

		sources = append(sources, &SendTarget{
			DisplayName: authenticatedUser.Username + " (you)",
			Value:       "user:" + authenticatedUser.ID,
		})

		var targets []*SendTarget
		targets = append(targets, &SendTarget{
			DisplayName: group.Name + " (group)",
			Value:       "group:" + group.ID,
		})

		for _, m := range ms {
			targetName := m.User.Username
			if m.User.ID == authenticatedUser.ID {
				targetName = targetName + " (you)"
			}
			targets = append(targets, &SendTarget{
				DisplayName: targetName,
				Value:       "user:" + m.UserID,
			})
			if m.UserID == authenticatedUser.ID {
				if m.IsAdmin() {
					sources = append(sources, &SendTarget{
						DisplayName: group.Name + " (group)",
						Value:       "group:" + m.GroupID,
					})
				}
			}
		}

		return c.Render(http.StatusOK, "group_send", map[string]interface{}{
			"Title":   "Hello",
			"Sources": sources,
			"Targets": targets,
		})
	}

	var payload SendOption
	if err := c.Bind(&payload); err != nil {
		return err
	}

	source, err := h.getTarget(payload.Source)
	if err != nil {
		return err
	}
	target, err := h.getTarget(payload.Target)
	if err != nil {
		return err
	}

	if payload.Type == Credits {

		amount, err := time.ParseDuration(payload.Amount)
		if err != nil {
			return err
		}

		credits := &api.Credits{
			ID:      uuid.NewV4().String(),
			GroupID: group.ID,
			SentTo:  target,
			SentBy:  source,
			Amount:  amount,
			Notes:   payload.Notes,
		}

		if err := h.creditsStore.Create(credits); err != nil {
			return err
		}

		if err := h.alertManager.AddAlert(c.Request(), c.Response().Writer, utils.Alert{
			Class:   "alert-success",
			Message: fmt.Sprintf("Successfully sent %s credits to %s", credits.Amount.String(), target.HTMLLink()),
		}); err != nil {
			return err
		}

	} else if payload.Type == Acknowledgement || payload.Type == Other {

		acknowledgementType := payload.AcknowledgementType
		if payload.Type == Other {
			acknowledgementType = api.Other
		}

		acknowledgement := &api.Acknowledgement{
			ID:      uuid.NewV4().String(),
			GroupID: group.ID,
			SentTo:  target,
			SentBy:  source,
			Type:    acknowledgementType,
			Notes:   payload.Notes,
		}

		if err := h.acknowledgementStore.Save(acknowledgement); err != nil {
			return err
		}

		if err := h.alertManager.AddAlert(c.Request(), c.Response().Writer, utils.Alert{
			Class:   "alert-success",
			Message: fmt.Sprintf("Successfully sent acknowledgement to %s", target.HTMLLink()),
		}); err != nil {
			return err
		}

	}

	c.Response().Header().Set("Location", fmt.Sprintf("%s://%s/groups/%s", c.Scheme(), c.Request().Host, group.ID))
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil

}
