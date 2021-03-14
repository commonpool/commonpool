package keys

import (
	"github.com/commonpool/backend/pkg/keys"
	domain2 "github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/trading/domain"
	"time"
)

type OfferReadModelBase struct {
	OfferKey            keys.OfferKey      `gorm:"primaryKey;type:varchar(128)" json:"offerId,omitempty"`
	GroupKey            keys.GroupKey      `gorm:"type:varchar(128)" json:"groupId"`
	Status              domain.OfferStatus `gorm:"type:varchar(128)" json:"status,omitempty"`
	Version             int                `json:"version,omitempty"`
	DeclinedAt          *time.Time         `json:"declinedAt,omitempty"`
	SubmittedAt         time.Time          `json:"submittedAt,omitempty"`
	ApprovedAt          *time.Time         `json:"approvedAt,omitempty"`
	CompletedAt         *time.Time         `json:"completedAt,omitempty"`
	StatusLastChangedAt *time.Time         `json:"last_status_changed_at;omitempty"`
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
	DeclinedBy  *OfferUserReadModel    `json:"declinedBy,omitempty"`
	SubmittedBy *OfferUserReadModel    `json:"submittedBy,omitempty"`
	OfferItems  []*OfferItemReadModel2 `json:"items"`
}

type OfferUserReadModel struct {
	UserKey  keys.UserKey `gorm:"primaryKey;type:varchar(128)" json:"userId"`
	Username string       `gorm:"type:varchar(128)" json:"username"`
	Version  int          `json:"version"`
}

type OfferResourceReadModel struct {
	ResourceKey  keys.ResourceKey     `gorm:"primaryKey;type:varchar(128)" json:"resourceId"`
	ResourceName string               `gorm:"type:varchar(128)" json:"resourceName"`
	Version      int                  `json:"version"`
	ResourceType domain2.ResourceType `gorm:"type:varchar(128)" json:"resourceType"`
	CallType     domain2.CallType     `gorm:"type:varchar(128)" json:"callType"`
	Owner        keys.Target          `json:"owner" gorm:"embedded;embeddedPrefix:owner_"`
}

type OfferGroupReadModel struct {
	GroupKey  keys.GroupKey `gorm:"primaryKey;type:varchar(128)" json:"groupId"`
	GroupName string        `gorm:"type:varchar(128)" json:"groupName"`
	Version   int           `json:"version"`
}

type OfferItemTargetReadModel struct {
	keys.Target
	GroupName    *string `gorm:"type:varchar(128)" json:"groupName,omitempty"`
	UserName     *string `gorm:"type:varchar(128)" json:"userName,omitempty"`
	UserVersion  *int    `json:"userVersion,omitempty"`
	GroupVersion *int    `json:"groupVersion,omitempty"`
	Name         string  `gorm:"type:varchar(128)" json:"name"`
}

func (o OfferItemTargetReadModel) GetName() string {
	if o.Type == keys.UserTarget && o.UserName != nil {
		return *o.UserName
	} else if o.Type == keys.GroupTarget && o.GroupName != nil {
		return *o.GroupName
	} else {
		return ""
	}
}

type OfferUserMembershipReadModel struct {
	UserKey  keys.UserKey  `gorm:"primaryKey;type:varchar(128)" json:"userId"`
	GroupKey keys.GroupKey `gorm:"primaryKey;type:varchar(128)" json:"groupId"`
	IsMember bool          `json:"is_member"`
	IsAdmin  bool          `json:"is_admin"`
	IsOwner  bool          `json:"is_owner"`
	Version  int           `json:"version"`
}

type OfferItemReadModelBase struct {
	OfferItemKey           keys.OfferItemKey    `gorm:"primaryKey;type:varchar(128)" json:"offerItemId,omitempty"`
	OfferKey               keys.OfferKey        `gorm:"not null;type:varchar(128)" json:"offerId,omitempty"`
	Version                int                  `json:"version,omitempty"`
	Type                   domain.OfferItemType `gorm:"not null;type:varchar(128)" json:"type,omitempty"`
	Amount                 *time.Duration       `json:"amount,omitempty"`
	Duration               *time.Duration       `json:"duration,omitempty"`
	ApprovedInbound        bool                 `gorm:"not null" json:"approvedInbound,omitempty"`
	ApprovedInboundAt      *time.Time           `json:"approvedInboundAt,omitempty"`
	ApprovedOutbound       bool                 `gorm:"not null" json:"approvedOutbound,omitempty"`
	ApprovedOutboundAt     *time.Time           `json:"approvedOutboundAt,omitempty"`
	ServiceGiven           bool                 `gorm:"not null" json:"serviceGiven,omitempty"`
	ServiceGivenAt         *time.Time           `json:"serviceGivenAt,omitempty"`
	ServiceReceived        bool                 `gorm:"not null" json:"serviceReceived,omitempty"`
	ServiceReceivedAt      *time.Time           `json:"serviceReceivedAt,omitempty"`
	ResourceGiven          bool                 `gorm:"not null" json:"resourceGiven,omitempty"`
	ResourceGivenAt        *time.Time           `json:"resourceGivenAt,omitempty"`
	ResourceTaken          bool                 `gorm:"not null" json:"resourceTaken,omitempty"`
	ResourceTakenAt        *time.Time           `json:"resourceTakenAt,omitempty"`
	ResourceBorrowed       bool                 `gorm:"not null" json:"resourceBorrowed,omitempty"`
	ResourceBorrowedAt     *time.Time           `json:"resourceBorrowedAt,omitempty"`
	ResourceLent           bool                 `gorm:"not null" json:"resourceLent,omitempty"`
	ResourceLentAt         *time.Time           `json:"resourceLentAt,omitempty"`
	BorrowedItemReturned   bool                 `gorm:"not null" json:"borrowedItemReturned,omitempty"`
	BorrowedItemReturnedAt *time.Time           `json:"borrowedItemReturnedAT,omitempty"`
	LentItemReceived       bool                 `gorm:"not null" json:"lentItemReceived,omitempty"`
	LentItemReceivedAt     *time.Time           `json:"lentItemReceivedAt,omitempty"`
}

type OfferItemReadModel2 struct {
	OfferItemReadModelBase
	ApprovedInboundBy      *OfferUserReadModel       `json:"approvedInboundBy,omitempty"`
	ApprovedOutboundBy     *OfferUserReadModel       `json:"approvedOutboundBy,omitempty"`
	ServiceGivenBy         *OfferUserReadModel       `json:"serviceGivenBy,omitempty"`
	ServiceReceivedBy      *OfferUserReadModel       `json:"serviceReceivedBy,omitempty"`
	ResourceGivenBy        *OfferUserReadModel       `json:"resourceGivenBy,omitempty"`
	ResourceTakenBy        *OfferUserReadModel       `json:"resourceTakenBy,omitempty"`
	ResourceBorrowedBy     *OfferUserReadModel       `json:"resourceBorrowedBy,omitempty"`
	ResourceLentBy         *OfferUserReadModel       `json:"resourceLentBy,omitempty"`
	BorrowedItemReturnedBy *OfferUserReadModel       `json:"borrowedItemReturnedBy,omitempty"`
	LentItemReceivedBy     *OfferUserReadModel       `json:"lentItemReceivedBy,omitempty"`
	From                   *OfferItemTargetReadModel `json:"from,omitempty"`
	To                     *OfferItemTargetReadModel `json:"to,omitempty"`
	Resource               *OfferResourceReadModel   `json:"resource,omitempty"`
}

func (o OfferItemReadModel2) GetType() domain.OfferItemType {
	return o.Type
}
func (o OfferItemReadModel2) GetOfferKey() keys.OfferKey {
	return o.OfferKey
}
func (o OfferItemReadModel2) GetKey() keys.OfferItemKey {
	return o.OfferItemKey
}

type OfferActionReadModel struct {
	Name         string             `json:"name"`
	Enabled      bool               `json:"enabled"`
	Completed    bool               `json:"completed"`
	ActionURL    string             `json:"actionUrl"`
	OfferItemKey *keys.OfferItemKey `json:"offerItemId"`
	Style        string             `json:"style"`
}

type OfferActionReadModels []OfferActionReadModel

type OfferReadModelWithActions struct {
	*OfferReadModel
	Actions OfferActionReadModels `json:"actions"`
}

type OfferItemReadModel struct {
	OfferItemReadModelBase
	ApprovedInboundBy      *keys.UserKey     `gorm:"type:varchar(128)"`
	ApprovedOutboundBy     *keys.UserKey     `gorm:"type:varchar(128)"`
	ServiceGivenBy         *keys.UserKey     `gorm:"type:varchar(128)"`
	ServiceReceivedBy      *keys.UserKey     `gorm:"type:varchar(128)"`
	ResourceGivenBy        *keys.UserKey     `gorm:"type:varchar(128)"`
	ResourceTakenBy        *keys.UserKey     `gorm:"type:varchar(128)"`
	ResourceBorrowedBy     *keys.UserKey     `gorm:"type:varchar(128)"`
	ResourceLentBy         *keys.UserKey     `gorm:"type:varchar(128)"`
	BorrowedItemReturnedBy *keys.UserKey     `gorm:"type:varchar(128)"`
	LentItemReceivedBy     *keys.UserKey     `gorm:"type:varchar(128)"`
	From                   *keys.Target      `gorm:"embedded;embeddedPrefix:from_"`
	To                     *keys.Target      `gorm:"embedded;embeddedPrefix:to_"`
	ResourceKey            *keys.ResourceKey `gorm:"type:varchar(128)"`
	ResourceName           string            `gorm:"type:varchar(128)"`
	ResourceVersion        int
}

type GroupReportItem struct {
	ID               string        `json:"id"`
	GroupKey         keys.GroupKey `json:"groupId"`
	Activity         string        `json:"activity"`
	GroupingID       string        `json:"groupingId"`
	ItemsReceived    int           `json:"itemsReceived"`
	ItemsGiven       int           `json:"itemsGiven"`
	ItemsOwned       int           `json:"itemsOwned"`
	ItemsLent        int           `json:"itemsLent"`
	ItemsBorrowed    int           `json:"itemsBorrowed"`
	ServicesGiven    int           `json:"servicesGiven"`
	ServicesReceived int           `json:"servicesReceived"`
	OfferCount       int           `json:"offerCount"`
	RequestsCount    int           `json:"requestCount"`
	HoursInBank      time.Duration `json:"hoursInBank"`
	EventTime        time.Time     `json:"eventTime"`
}
