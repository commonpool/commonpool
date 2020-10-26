package handler

import (
	"encoding/json"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateGroup(t *testing.T) {

	mockLoggedInAs(user1)
	createGroupResponse := createGroup(t, "name", "description")
	groupId := createGroupResponse.Group.ID

	getGroupResponse := getGroup(t, groupId)
	getMembershipsResponse := getLoggedInUserMemberships(t)

	assert.Equal(t, createGroupResponse.Group.Name, "name")
	assert.Equal(t, createGroupResponse.Group.Description, "description")
	assert.Equal(t, getGroupResponse.Group.Name, "name")
	assert.Equal(t, getGroupResponse.Group.Description, "description")
	assert.Equal(t, 1, len(getMembershipsResponse.Memberships))
	assert.Equal(t, true, getMembershipsResponse.Memberships[0].UserConfirmed)
	assert.Equal(t, true, getMembershipsResponse.Memberships[0].GroupConfirmed)
	assert.Equal(t, false, getMembershipsResponse.Memberships[0].IsDeactivated)
	assert.Equal(t, user1.Subject, getMembershipsResponse.Memberships[0].UserID)
	assert.Equal(t, true, getMembershipsResponse.Memberships[0].IsMember)
	assert.Equal(t, true, getMembershipsResponse.Memberships[0].IsAdmin)
	assert.Equal(t, groupId, getMembershipsResponse.Memberships[0].GroupID)
	assert.Equal(t, "name", getMembershipsResponse.Memberships[0].GroupName)

}

func createGroup(t *testing.T, name string, description string) web.CreateGroupResponse {
	request := web.CreateGroupRequest{
		Name:        name,
		Description: description,
	}
	js, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	jsStr := string(js)
	_, _, rec, c := newRequest(echo.POST, "/api/v1/groups", &jsStr)
	err = h.CreateGroup(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	response := web.CreateGroupResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	return response
}

func getGroup(t *testing.T, id string) web.GetGroupResponse {
	_, _, rec, c := newRequest(echo.POST, "/api/v1/groups/"+id, nil)
	c.SetParamNames("id")
	c.SetParamValues(id)
	err := h.GetGroup(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	response := web.GetGroupResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	return response
}

func getLoggedInUserMemberships(t *testing.T) web.GetUserMembershipsResponse {
	_, _, rec, c := newRequest(echo.POST, "/api/v1/my/memberships", nil)
	err := h.GetLoggedInUserMemberships(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	response := web.GetUserMembershipsResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	return response
}