package listeners

import (
	"context"
	"database/sql"
	userdomain "github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	groupdomain "github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	resourcedomain "github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/trading/domain"
	groupreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"time"
)

type OfferReadModelHandler struct {
	catchUpFactory eventbus.CatchUpListenerFactory
	db             *gorm.DB
}

func NewOfferReadModelHandler(db *gorm.DB, catchUpFactory eventbus.CatchUpListenerFactory) *OfferReadModelHandler {
	var transactionHistoryHandler = &OfferReadModelHandler{
		catchUpFactory: catchUpFactory,
		db:             db,
	}
	return transactionHistoryHandler
}

func (h *OfferReadModelHandler) Start(ctx context.Context) error {
	err := h.migrateDatabase()
	if err != nil {
		return err
	}
	listener := h.catchUpFactory("offer-read-model", time.Second*10)
	eventTypes := domain.AllEvents
	eventTypes = append(eventTypes, resourcedomain.ResourceRegisteredEvent)
	eventTypes = append(eventTypes, resourcedomain.ResourceInfoChangedEvent)
	eventTypes = append(eventTypes, userdomain.UserDiscoveredEvent)
	eventTypes = append(eventTypes, userdomain.UserInfoChangedEvent)
	eventTypes = append(eventTypes, groupdomain.GroupCreatedEvent)
	eventTypes = append(eventTypes, groupdomain.GroupInfoChangedEvent)
	eventTypes = append(eventTypes, groupdomain.MembershipStatusChangedEvent)
	if err := listener.Initialize(ctx, "offer-read-model", eventTypes); err != nil {
		return err
	}
	return listener.Listen(ctx, h.HandleEvents)
}

func (h *OfferReadModelHandler) migrateDatabase() error {
	if err := h.db.AutoMigrate(
		&groupreadmodels.OfferUserReadModel{},
		&groupreadmodels.OfferResourceReadModel{},
		&groupreadmodels.DBOfferReadModel{},
		&groupreadmodels.OfferItemReadModel{},
		&groupreadmodels.OfferGroupReadModel{},
		&groupreadmodels.OfferUserMembershipReadModel{},
	); err != nil {
		return err
	}
	return nil
}

func (h *OfferReadModelHandler) HandleEvents(ctx context.Context, events []eventsource.Event) error {
	for _, event := range events {
		if err := h.HandleEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (h *OfferReadModelHandler) HandleEvent(event eventsource.Event) error {
	switch e := event.(type) {
	case *domain.OfferSubmitted:
		return h.handleOfferSubmitted(e)
	case *domain.OfferApproved:
		return h.handleOfferApproved(e)
	case *domain.OfferDeclined:
		return h.handleOfferDeclined(e)
	case *domain.OfferCompleted:
		return h.handleOfferCompleted(e)
	case *domain.OfferItemApproved:
		return h.handleOfferItemApproved(e)
	case *domain.ResourceTransferGivenNotified:
		return h.handleResourceTransferGivenNotified(e)
	case *domain.ResourceTransferReceivedNotified:
		return h.handleResourceTransferReceivedNotified(e)
	case *domain.ServiceGivenNotified:
		return h.handleServiceGivenNotified(e)
	case *domain.ServiceReceivedNotified:
		return h.handleServiceReceivedNotified(e)
	case *domain.ResourceBorrowedNotified:
		return h.handleResourceBorrowedNotified(e)
	case *domain.ResourceLentNotified:
		return h.handleResourceLentNotified(e)
	case *domain.BorrowerReturnedResourceNotified:
		return h.handleBorrowerReturnedResourceNotified(e)
	case *domain.LenderReceivedBackResourceNotified:
		return h.handleLenderReceivedBackResourceNotified(e)
	case userdomain.UserDiscovered:
		return h.handleUserDiscovered(e)
	case userdomain.UserInfoChanged:
		return h.handleUserInfoChanged(e)
	case resourcedomain.ResourceRegistered:
		return h.handleResourceRegistered(e)
	case resourcedomain.ResourceInfoChanged:
		return h.handleResourceInfoChanged(e)
	case groupdomain.GroupCreated:
		return h.handleGroupCreated(e)
	case groupdomain.GroupInfoChanged:
		return h.handleGroupInfoChanged(e)
	case groupdomain.MembershipStatusChanged:
		return h.handleMembershipStatusChanged(e)
	default:
		log.Warnf("unhandled event type: %s", reflect.TypeOf(e).String())
	}
	return nil
}

func (h *OfferReadModelHandler) handleLenderReceivedBackResourceNotified(e *domain.LenderReceivedBackResourceNotified) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), map[string]interface{}{
		"lent_item_received":    true,
		"lent_item_received_by": e.NotifiedBy,
		"lent_item_received_at": e.GetEventTime(),
	})
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleBorrowerReturnedResourceNotified(e *domain.BorrowerReturnedResourceNotified) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), map[string]interface{}{
		"borrowed_item_returned":    true,
		"borrowed_item_returned_by": e.NotifiedBy,
		"borrowed_item_returned_at": e.GetEventTime(),
	})
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleResourceLentNotified(e *domain.ResourceLentNotified) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), map[string]interface{}{
		"resource_lent":    true,
		"resource_lent_by": e.NotifiedBy,
		"resource_lent_at": e.GetEventTime(),
	})
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleResourceBorrowedNotified(e *domain.ResourceBorrowedNotified) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), map[string]interface{}{
		"resource_borrowed":    true,
		"resource_borrowed_by": e.NotifiedBy,
		"resource_borrowed_at": e.GetEventTime(),
	})
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleServiceReceivedNotified(e *domain.ServiceReceivedNotified) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), map[string]interface{}{
		"service_received":    true,
		"service_received_by": e.NotifiedBy,
		"service_received_at": e.GetEventTime(),
	})
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleServiceGivenNotified(e *domain.ServiceGivenNotified) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), map[string]interface{}{
		"service_given":    true,
		"service_given_by": e.ServiceGivenNotifiedPayload.NotifiedBy,
		"service_given_at": e.GetEventTime(),
	})
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleResourceTransferReceivedNotified(e *domain.ResourceTransferReceivedNotified) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), map[string]interface{}{
		"resource_taken":    true,
		"resource_taken_by": e.NotifiedBy,
		"resource_taken_at": e.GetEventTime(),
	})
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleResourceTransferGivenNotified(e *domain.ResourceTransferGivenNotified) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), map[string]interface{}{
		"resource_given":    true,
		"resource_given_by": e.NotifiedBy,
		"resource_given_at": e.GetEventTime(),
	})
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleOfferItemApproved(e *domain.OfferItemApproved) error {
	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"version": e.GetSequenceNo(),
	}
	if e.Direction == domain.Inbound {
		updates["approved_inbound"] = true
		updates["approved_inbound_by"] = e.ApprovedBy.String()
		updates["approved_inbound_at"] = e.GetEventTime()
	} else if e.Direction == domain.Outbound {
		updates["approved_outbound"] = true
		updates["approved_outbound_by"] = e.ApprovedBy.String()
		updates["approved_outbound_at"] = e.GetEventTime()
	}
	err = h.updateOfferItem(e.OfferItemKey.String(), e.GetSequenceNo(), updates)
	if err != nil {
		return err
	}
	return h.updateOfferVersion(h.db, offerKey, e.SequenceNo)
}

func (h *OfferReadModelHandler) handleOfferCompleted(e *domain.OfferCompleted) error {
	return h.db.Model(&groupreadmodels.DBOfferReadModel{}).
		Where("offer_key = ? and version < ?", e.GetAggregateID(), e.GetSequenceNo()).
		Updates(map[string]interface{}{
			"status":       "completed",
			"version":      e.GetSequenceNo(),
			"completed_at": e.GetEventTime(),
		}).Error
}

func (h *OfferReadModelHandler) handleOfferApproved(e *domain.OfferApproved) error {
	return h.db.Model(&groupreadmodels.DBOfferReadModel{}).
		Where("offer_key = ? and version < ?", e.GetAggregateID(), e.GetSequenceNo()).
		Updates(map[string]interface{}{
			"status":      "approved",
			"version":     e.GetSequenceNo(),
			"approved_at": e.GetEventTime(),
		}).Error
}

func (h *OfferReadModelHandler) updateOfferVersion(db *gorm.DB, offerKey keys.OfferKey, version int) error {
	return db.Model(&groupreadmodels.DBOfferReadModel{}).Where("offer_key = ? and version < ?", offerKey, version).Update("version", version).Error
}

func (h *OfferReadModelHandler) handleOfferDeclined(e *domain.OfferDeclined) error {
	return h.db.Model(&groupreadmodels.DBOfferReadModel{}).
		Where("offer_key = ? and version < ?", e.GetAggregateID(), e.GetSequenceNo()).
		Updates(map[string]interface{}{
			"status":      "declined",
			"version":     e.GetSequenceNo(),
			"declined_by": e.OfferDeclinedPayload.DeclinedBy,
			"declined_at": e.GetEventTime(),
		}).Error
}

func (h *OfferReadModelHandler) updateOfferItem(offerItemId string, expectedVersion int, updates map[string]interface{}) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		updates["version"] = expectedVersion
		return tx.
			Model(&groupreadmodels.OfferItemReadModel{}).
			Where("offer_item_key = ? and version < ?", offerItemId, expectedVersion).
			Updates(updates).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (h *OfferReadModelHandler) handleOfferSubmitted(e *domain.OfferSubmitted) error {

	offerKey, err := keys.ParseOfferKey(e.AggregateID)
	if err != nil {
		return err
	}

	return h.db.Transaction(func(tx *gorm.DB) error {

		getOptimisticLocking(tx, e.GetSequenceNo(), []clause.Column{
			{Name: "offer_key"},
		}).Create(&groupreadmodels.DBOfferReadModel{
			OfferReadModelBase: groupreadmodels.OfferReadModelBase{
				OfferKey:    offerKey,
				GroupKey:    e.GroupKey,
				Status:      domain.Pending,
				Version:     e.GetSequenceNo(),
				SubmittedAt: e.GetEventTime(),
			},
			SubmittedBy: &e.SubmittedBy,
		})

		var offerItems []*groupreadmodels.OfferItemReadModel

		for _, offerItem := range e.OfferItems.Items {

			rm := &groupreadmodels.OfferItemReadModel{
				OfferItemReadModelBase: groupreadmodels.OfferItemReadModelBase{
					OfferItemKey:     offerItem.GetKey(),
					OfferKey:         offerKey,
					Version:          e.GetSequenceNo(),
					Type:             offerItem.Type(),
					ApprovedInbound:  false,
					ApprovedOutbound: false,
				},
			}

			if resourceTransfer, ok := offerItem.AsResourceTransfer(); ok {
				rm.To = resourceTransfer.To
				rm.ResourceKey = &resourceTransfer.ResourceKey
			}

			if provideService, ok := offerItem.AsProvideService(); ok {
				rm.To = provideService.To
				rm.ResourceKey = &provideService.ResourceKey
				rm.Duration = &provideService.Duration
			}

			if borrowResource, ok := offerItem.AsBorrowResource(); ok {
				rm.To = borrowResource.To
				rm.ResourceKey = &borrowResource.ResourceKey
				rm.Duration = &borrowResource.Duration
			}

			if creditTransfer, ok := offerItem.AsCreditTransfer(); ok {
				rm.To = creditTransfer.To
				rm.From = creditTransfer.From
				rm.Amount = &creditTransfer.Amount
			}

			offerItems = append(offerItems, rm)
		}

		if err := getOptimisticLocking(tx, e.GetSequenceNo(), []clause.Column{
			{Name: "offer_item_key"},
		}).Create(offerItems).Error; err != nil {
			return err
		}

		return nil

	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (h *OfferReadModelHandler) handleUserDiscovered(e userdomain.UserDiscovered) error {
	return h.createUser(e.AggregateID, e.UserInfo.Username, e.SequenceNo)
}

func (h *OfferReadModelHandler) createUser(id string, username string, version int) error {
	userKey := keys.NewUserKey(id)
	return getOptimisticLocking(h.db, version, []clause.Column{
		{Name: "user_key"},
	}).
		Model(&groupreadmodels.OfferUserReadModel{}).
		Create(&groupreadmodels.OfferUserReadModel{
			UserKey:  userKey,
			Username: username,
			Version:  version,
		}).Error
}

func (h *OfferReadModelHandler) handleUserInfoChanged(e userdomain.UserInfoChanged) error {
	if e.OldUserInfo.Username != e.NewUserInfo.Username {
		userKey := keys.NewUserKey(e.AggregateID)
		qry := h.db.Model(&groupreadmodels.OfferUserReadModel{}).
			Where("version < ? and user_key = ?", e.SequenceNo, userKey).
			Updates(map[string]interface{}{
				"username": e.NewUserInfo.Username,
				"version":  e.SequenceNo,
			})
		if qry.Error != nil {
			return qry.Error
		}
	}
	return nil
}

func (h *OfferReadModelHandler) handleResourceRegistered(e resourcedomain.ResourceRegistered) error {
	resourceKey, err := keys.ParseResourceKey(e.AggregateID)
	if err != nil {
		return err
	}
	return getOptimisticLocking(h.db, e.SequenceNo, []clause.Column{
		{Name: "resource_key"},
	}).
		Model(&groupreadmodels.OfferResourceReadModel{}).
		Create(&groupreadmodels.OfferResourceReadModel{
			ResourceKey:  resourceKey,
			ResourceName: e.ResourceInfo.Name,
			Version:      e.SequenceNo,
			ResourceType: e.ResourceInfo.ResourceType,
			CallType:     e.ResourceInfo.CallType,
			Owner:        e.RegisteredFor,
		}).Error
}

func (h *OfferReadModelHandler) handleResourceInfoChanged(e resourcedomain.ResourceInfoChanged) error {
	resourceKey, err := keys.ParseResourceKey(e.AggregateID)
	if err != nil {
		return err
	}
	if e.OldResourceInfo.Name != e.NewResourceInfo.Name {
		qry := h.db.Model(&groupreadmodels.OfferResourceReadModel{}).
			Where("resource_key = ? and version < ?", resourceKey, e.SequenceNo).
			Updates(map[string]interface{}{
				"resource_name": e.NewResourceInfo.Name,
				"version":       e.SequenceNo,
			})
		if qry.Error != nil {
			return qry.Error
		}
	}
	return nil
}

func (h *OfferReadModelHandler) handleGroupInfoChanged(e groupdomain.GroupInfoChanged) error {
	groupKey, err := keys.ParseGroupKey(e.AggregateID)
	if err != nil {
		return err
	}
	if e.OldGroupInfo.Name != e.NewGroupInfo.Name {
		qry := h.db.Model(&groupreadmodels.OfferGroupReadModel{}).
			Where("group_key = ? and version < ?", groupKey, e.SequenceNo).
			Updates(map[string]interface{}{
				"group_name": e.NewGroupInfo.Name,
				"version":    e.SequenceNo,
			})
		if qry.Error != nil {
			return qry.Error
		}
	}
	return nil
}

func (h *OfferReadModelHandler) handleGroupCreated(e groupdomain.GroupCreated) error {
	groupKey, err := keys.ParseGroupKey(e.AggregateID)
	if err != nil {
		return err
	}
	return getOptimisticLocking(h.db, e.SequenceNo, []clause.Column{
		{Name: "resource_key"},
	}).
		Model(&groupreadmodels.OfferGroupReadModel{}).
		Create(&groupreadmodels.OfferGroupReadModel{
			GroupKey:  groupKey,
			GroupName: e.GroupInfo.Name,
			Version:   e.SequenceNo,
		}).Error
}

func (h *OfferReadModelHandler) handleMembershipStatusChanged(e groupdomain.MembershipStatusChanged) error {
	memberKey := e.MemberKey
	groupKey, err := keys.ParseGroupKey(e.AggregateID)
	if err != nil {
		return err
	}
	if e.IsNewMembership {
		return getOptimisticLocking(h.db, e.SequenceNo, []clause.Column{}).
			Model(&groupreadmodels.OfferUserMembershipReadModel{}).
			Create(&groupreadmodels.OfferUserMembershipReadModel{
				UserKey:  memberKey,
				GroupKey: groupKey,
				IsMember: e.NewPermissions.IsMember(),
				IsAdmin:  e.NewPermissions.IsAdmin(),
				IsOwner:  e.NewPermissions.IsOwner(),
				Version:  e.SequenceNo,
			}).Error
	} else if e.IsCanceledMembership {
		return h.db.Model(&groupreadmodels.OfferUserMembershipReadModel{}).
			Delete(
				&groupreadmodels.OfferUserMembershipReadModel{},
				"user_key = ? and group_key = ? and version < ?",
				memberKey,
				groupKey,
				e.SequenceNo,
			).Error
	} else {
		return h.db.Model(&groupreadmodels.OfferUserMembershipReadModel{}).
			Where("user_key = ? and version < ?", e.MemberKey, e.SequenceNo).
			Updates(map[string]interface{}{
				"is_member": e.NewPermissions.IsMember(),
				"is_owner":  e.NewPermissions.IsOwner(),
				"is_admin":  e.NewPermissions.IsAdmin(),
				"version":   e.SequenceNo,
			}).Error
	}
}

func getOptimisticLocking(db *gorm.DB, version int, columns []clause.Column) *gorm.DB {
	return db.Clauses(
		clause.OnConflict{
			Where: clause.Where{
				Exprs: []clause.Expression{
					clause.Lt{
						Column: "version",
						Value:  version,
					},
				},
			},
			UpdateAll: true,
		})
}
