package chat

import (
	"github.com/commonpool/backend/errors"
	"net/http"
)

var ErrChannelNotFound = errors.NewWebServiceException("channel not found", "ErrChannelNotFound", http.StatusNotFound)
