package web

type GetTradingHistoryRequest struct {
	UserIDs []string `json:"userIds"`
}

type HistoryMatrix struct {
}

type HistoryEntry struct {
	Activity          string `json:"activity"`
	ItemsReceived     int    `json:"itemsReceived"`
	ItemsGiven        int    `json:"itemsGiven"`
	ServicesGiven     int    `json:"servicesGiven"`
	ServicesReceived  int    `json:"servicesReceived"`
	ItemsBorrowed     int    `json:"itemsBorrowed"`
	ItemsOwnedByGroup int    `json:"itemsOwnedByGroup"`
	OffersInGroup     int    `json:"offersInGroup"`
	HoursInBank       int    `json:"hoursInBank"`
}

type TradingHistoryEntry struct {
	Timestamp         string  `json:"timestamp"`
	FromUserID        string  `json:"fromUserId"`
	FromUsername      string  `json:"fromUsername"`
	ToUserID          string  `json:"toUserId"`
	ToUsername        string  `json:"toUsername"`
	ResourceID        *string `json:"resourceId"`
	TimeAmountSeconds *int64  `json:"timeAmountSeconds"`
}

type GetTradingHistoryResponse struct {
	Entries []TradingHistoryEntry `json:"entries"`
}
