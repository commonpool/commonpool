package trading

import (
	"github.com/commonpool/backend/errors"
	"net/http"
)

var ErrUserNotPartOfOfferItem = errors.NewWebServiceException("you can't confirm an item you're not receiving or giving", "ErrUserNotPartOfOfferItem", http.StatusForbidden)
