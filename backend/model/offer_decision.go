package model

import (
	uuid "github.com/satori/go.uuid"
)

type Decision int

const (
	PendingDecision Decision = iota
	AcceptedDecision
	DeclinedDecision
)

type OfferDecision struct {
	OfferID  uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID   string    `gorm:"primary_key"`
	Decision Decision
}

func (d *OfferDecision) GetKey() OfferDecisionKey {
	return NewOfferDecisionKey(d.GetOfferKey(), d.GetUserKey())
}

func (d *OfferDecision) GetOfferKey() OfferKey {
	return NewOfferKey(d.OfferID)
}
func (d *OfferDecision) GetUserKey() UserKey {
	return NewUserKey(d.UserID)
}
