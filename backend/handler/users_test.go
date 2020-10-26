package handler

import (
	"encoding/json"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUserPicker(t *testing.T) {

	_, _, rec, c := newRequest(echo.GET, "/api/v1/users/search", nil)
	c.QueryParams().Set("skip", "0")
	c.QueryParams().Set("take", "10")
	c.QueryParams().Set("query", "user1")

	err := h.SearchUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	res := web.UsersInfoResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, 1, len(res.Users))
	assert.Equal(t, 0, res.Skip)
	assert.Equal(t, 10, res.Take)

}
