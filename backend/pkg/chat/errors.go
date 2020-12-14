package chat

import (
	"github.com/commonpool/backend/pkg/exceptions"
	"net/http"
)

var ErrChannelNotFound = exceptions.NewWebServiceException("channel not found", "ErrChannelNotFound", http.StatusNotFound)
