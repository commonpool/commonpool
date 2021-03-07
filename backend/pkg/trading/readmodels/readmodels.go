package keys

import (
	"github.com/commonpool/backend/pkg/keys"
	domain2 "github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/trading/domain"
	"time"
)

type OfferReadModelBase struct {
	OfferKey    keys.OfferKey      `gorm:"primaryKey" json:"offer_key,omitempty"`
	Status      domain.OfferStatus `json:"status,omitempty"`
	Version     int                `json:"version,omitempty"`
	DeclinedAt  *time.Time         `json:"declined_at,omitempty"`
	SubmittedAt time.Time          `json:"submitted_at,omitempty"`
	ApprovedAt  *time.Time         `json:"approved_at,omitempty"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
}

type DBOfferReadModel struct {
	OfferReadModelBase
	DeclinedBy  *keys.UserKey
	SubmittedBy *keys.UserKey
}

func (DBOfferReadModel) TableName() string {
	return "offer_read_models"
}

type OfferReadModel struct {
	OfferReadModelBase
	DeclinedBy  *OfferUserReadModel    `json:"declined_by,omitempty"`
	SubmittedBy *OfferUserReadModel    `json:"submitted_by,omitempty"`
	OfferItems  []*OfferItemReadModel2 `json:"offer_items"`
}

type OfferUserReadModel struct {
	UserKey  keys.UserKey `gorm:"primaryKey" json:"user_key"`
	Username string       `json:"username"`
	Version  int          `json:"version"`
}

type OfferResourceReadModel struct {
	ResourceKey  keys.ResourceKey     `gorm:"primaryKey" json:"resource_key"`
	ResourceName string               `json:"resource_name"`
	Version      int                  `json:"version"`
	ResourceType domain2.ResourceType `json:"resource_type"`
	CallType     domain2.CallType     `json:"call_type"`
	Owner        domain.Target        `json:"owner" gorm:"embedded;embeddedPrefix:owner_"`
}

type OfferGroupReadModel struct {
	GroupKey  keys.GroupKey `gorm:"primaryKey" json:"group_key"`
	GroupName string        `json:"group_name"`
	Version   int           `json:"version"`
}

type OfferItemTargetReadModel struct {
	domain.Target
	GroupName    *string `json:"group_name,omitempty"`
	UserName     *string `json:"user_name,omitempty"`
	UserVersion  *int    `json:"user_version,omitempty"`
	GroupVersion *int    `json:"group_version,omitempty"`
}

type OfferUserMembershipReadModel struct {
	UserKey  keys.UserKey  `gorm:"primaryKey" json:"user_key"`
	GroupKey keys.GroupKey `gorm:"primaryKey" json:"group_key"`
	IsMember bool          `json:"is_member"`
	IsAdmin  bool          `json:"is_admin"`
	IsOwner  bool          `json:"is_owner"`
	Version  int           `json:"version"`
}

type OfferItemReadModelBase struct {
	OfferItemKey           keys.OfferItemKey    `gorm:"primaryKey" json:"offer_item_key,omitempty"`
	OfferKey               keys.OfferKey        `gorm:"not null" json:"offer_key,omitempty"`
	Version                int                  `json:"version,omitempty"`
	Type                   domain.OfferItemType `gorm:"not null" json:"type,omitempty"`
	Amount                 *time.Duration       `json:"amount,omitempty"`
	Duration               *time.Duration       `json:"duration,omitempty"`
	ApprovedInbound        bool                 `gorm:"not null" json:"approved_inbound,omitempty"`
	ApprovedInboundAt      *time.Time           `json:"approved_inbound_at,omitempty"`
	ApprovedOutbound       bool                 `gorm:"not null" json:"approved_outbound,omitempty"`
	ApprovedOutboundAt     *time.Time           `json:"approved_outbound_at,omitempty"`
	ServiceGiven           bool                 `gorm:"not null" json:"service_given,omitempty"`
	ServiceGivenAt         *time.Time           `json:"service_given_at,omitempty"`
	ServiceReceived        bool                 `gorm:"not null" json:"service_received,omitempty"`
	ServiceReceivedAt      *time.Time           `json:"service_received_at,omitempty"`
	ResourceGiven          bool                 `gorm:"not null" json:"resource_given,omitempty"`
	ResourceGivenAt        *time.Time           `json:"resource_given_at,omitempty"`
	ResourceTaken          bool                 `gorm:"not null" json:"resource_taken,omitempty"`
	ResourceTakenAt        *time.Time           `json:"resource_taken_at,omitempty"`
	ResourceBorrowed       bool                 `gorm:"not null" json:"resource_borrowed,omitempty"`
	ResourceBorrowedAt     *time.Time           `json:"resource_borrowed_at,omitempty"`
	ResourceLent           bool                 `gorm:"not null" json:"resource_lent,omitempty"`
	ResourceLentAt         *time.Time           `json:"resource_lent_at,omitempty"`
	BorrowedItemReturned   bool                 `gorm:"not null" json:"borrowed_item_returned,omitempty"`
	BorrowedItemReturnedAt *time.Time           `json:"borrowed_item_returned_at,omitempty"`
	LentItemReceived       bool                 `gorm:"not null" json:"lent_item_received,omitempty"`
	LentItemReceivedAt     *time.Time           `json:"lent_item_received_at,omitempty"`
}

type OfferItemReadModel2 struct {
	OfferItemReadModelBase
	ApprovedInboundBy      *OfferUserReadModel       `json:"approved_inbound_by,omitempty"`
	ApprovedOutboundBy     *OfferUserReadModel       `json:"approved_outbound_by,omitempty"`
	ServiceGivenBy         *OfferUserReadModel       `json:"service_given_by,omitempty"`
	ServiceReceivedBy      *OfferUserReadModel       `json:"service_received_by,omitempty"`
	ResourceGivenBy        *OfferUserReadModel       `json:"resource_given_by,omitempty"`
	ResourceTakenBy        *OfferUserReadModel       `json:"resource_taken_by,omitempty"`
	ResourceBorrowedBy     *OfferUserReadModel       `json:"resource_borrowed_by,omitempty"`
	ResourceLentBy         *OfferUserReadModel       `json:"resource_lent_by,omitempty"`
	BorrowedItemReturnedBy *OfferUserReadModel       `json:"borrowed_item_returned_by,omitempty"`
	LentItemReceivedBy     *OfferUserReadModel       `json:"lent_item_received_by,omitempty"`
	From                   *OfferItemTargetReadModel `json:"from,omitempty"`
	To                     *OfferItemTargetReadModel `json:"to,omitempty"`
	Resource               *OfferResourceReadModel   `json:"resource"`
}

type OfferItemReadModel struct {
	OfferItemReadModelBase
	ApprovedInboundBy      *keys.UserKey
	ApprovedOutboundBy     *keys.UserKey
	ServiceGivenBy         *keys.UserKey
	ServiceReceivedBy      *keys.UserKey
	ResourceGivenBy        *keys.UserKey
	ResourceTakenBy        *keys.UserKey
	ResourceBorrowedBy     *keys.UserKey
	ResourceLentBy         *keys.UserKey
	BorrowedItemReturnedBy *keys.UserKey
	LentItemReceivedBy     *keys.UserKey
	From                   *domain.Target `gorm:"embedded;embeddedPrefix:from_"`
	To                     *domain.Target `gorm:"embedded;embeddedPrefix:to_"`
	ResourceKey            *keys.ResourceKey
}
