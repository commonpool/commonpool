package chat

import (
	"github.com/commonpool/backend/model"
	"net/http"
)

var ErrChannelNotFound = model.NewWebServiceException("channel not found", "ErrChannelNotFound", http.StatusNotFound)
