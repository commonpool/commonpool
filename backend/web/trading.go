package web

type GetTradingHistoryRequest struct {
	UserIDs []string `json:"userIds"`
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
