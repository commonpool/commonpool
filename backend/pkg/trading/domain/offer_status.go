package domain

type OfferStatus string

const (
	Pending   OfferStatus = "pending"
	Declined  OfferStatus = "declined"
	Approved  OfferStatus = "approved"
	Completed OfferStatus = "completed"
)
