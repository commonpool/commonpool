package trading

type OfferStatus int

const (
	PendingOffer OfferStatus = iota
	AcceptedOffer
	CanceledOffer
	DeclinedOffer
	ExpiredOffer
	CompletedOffer
)
