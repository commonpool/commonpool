package listeners

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
	resourcedomain "github.com/commonpool/backend/pkg/resource/domain"
	resourcereadmodels "github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/commonpool/backend/pkg/trading/domain"
	readmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"math"
	"reflect"
	"strings"
	"time"
)

type GroupReportListener struct {
	catchUpFactory eventbus.CatchUpListenerFactory
	db             *gorm.DB
}

func NewGroupReportListener(
	catchUpFactory eventbus.CatchUpListenerFactory,
	db *gorm.DB) *GroupReportListener {
	return &GroupReportListener{
		catchUpFactory: catchUpFactory,
		db:             db,
	}
}

func (h *GroupReportListener) Start(ctx context.Context) error {
	err := h.migrateDatabase()
	if err != nil {
		return err
	}
	listener := h.catchUpFactory("group-report-readmodel", time.Second*10)
	var eventTypes []string
	eventTypes = append(eventTypes, domain.OfferCompletedEvent)
	eventTypes = append(eventTypes, resourcedomain.ResourceGroupSharingChangedEvent)
	if err := listener.Initialize(ctx, "group-report-readmodel", eventTypes); err != nil {
		return err
	}
	return listener.Listen(ctx, h.HandleEvents)
}

func (h *GroupReportListener) migrateDatabase() error {
	if err := h.db.AutoMigrate(
		&readmodels.GroupReportItem{},
	); err != nil {
		return err
	}
	if err := h.db.Delete(&readmodels.GroupReportItem{}, "1 = 1").Error; err != nil {
		return err
	}
	return nil
}

func (h *GroupReportListener) HandleEvents(ctx context.Context, events []eventsource.Event) error {
	for _, event := range events {
		if err := h.HandleEvent(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (h *GroupReportListener) HandleEvent(ctx context.Context, event eventsource.Event) error {
	switch e := event.(type) {
	case *domain.OfferCompleted:
		return h.handleOfferCompleted(ctx, e)
	case resourcedomain.ResourceGroupSharingChanged:
		return h.handleResourceGroupSharingChanged(ctx, e)
	default:
		log.Warnf("unhandled event type: %s", reflect.TypeOf(e).String())
	}
	return nil
}

func (h *GroupReportListener) handleOfferCompleted(ctx context.Context, e *domain.OfferCompleted) error {

	var resourceKeys []keys.ResourceKey
	for _, item := range e.OfferCompletedPayload.OfferItems.Items {
		if rkg, ok := item.(domain.ResourceKeyGetter); ok {
			resourceKeys = append(resourceKeys, rkg.GetResourceKey())
		}
	}

	var resourceMap = map[keys.ResourceKey]*resourcereadmodels.DbResourceReadModel{}
	if len(resourceKeys) > 0 {
		var sb strings.Builder
		sb.WriteString("resource_key in (")
		var params []interface{}
		for i, resourceKey := range resourceKeys {
			sb.WriteString("?")
			params = append(params, resourceKey)
			if i < len(resourceKeys)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
		var rrm []*resourcereadmodels.DbResourceReadModel
		err := h.db.
			Where(sb.String(), params...).
			Model(&resourcereadmodels.DbResourceReadModel{}).
			Find(&rrm).Error
		if err != nil {
			return err
		}
		for _, resourceReadModel := range rrm {
			resourceMap[resourceReadModel.ResourceKey] = resourceReadModel
		}
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, item := range e.OfferItems.Items {

		switch item.GetType() {
		case domain.CreditTransfer:
			i := item.(*domain.CreditTransferItem)
			if i.From.IsGroup() || i.To.IsGroup() {
				var credit time.Duration = 0
				var activity string
				hours := i.Amount.Truncate(time.Hour).Hours()
				hours = math.Round(hours*100) / 100
				if i.From.IsGroup() {
					credit = -i.Amount
					activity = fmt.Sprintf("group gave %.2fh. of timebank credits", hours)
				}
				if i.To.IsGroup() {
					credit = i.Amount
					activity = fmt.Sprintf("group received %.2fh. of timebank credits", hours)
				}
				g.Go(func() error {
					return h.db.Create(&readmodels.GroupReportItem{
						ID:          uuid.NewV4().String(),
						GroupKey:    e.GroupKey,
						Activity:    activity,
						HoursInBank: credit,
						GroupingID:  e.AggregateID,
						EventTime:   e.EventTime,
					}).Error
				})
			}
		case domain.ResourceTransfer:
			i := item.(*domain.ResourceTransferItem)
			if i.To.IsGroup() {
				g.Go(func() error {
					return h.db.Create(&readmodels.GroupReportItem{
						ID:            uuid.NewV4().String(),
						GroupKey:      e.GroupKey,
						Activity:      fmt.Sprintf("group received resource '%s'", resourceMap[i.ResourceKey].Name),
						GroupingID:    e.AggregateID,
						ItemsReceived: 1,
						EventTime:     e.EventTime,
					}).Error
				})
			}
			resource := resourceMap[i.ResourceKey]
			if resource.Owner.IsGroup() {
				g.Go(func() error {
					return h.db.Create(&readmodels.GroupReportItem{
						ID:         uuid.NewV4().String(),
						GroupKey:   resource.Owner.GetGroupKey(),
						Activity:   fmt.Sprintf("group gave resource '%s'", resourceMap[i.ResourceKey].Name),
						GroupingID: e.AggregateID,
						ItemsGiven: 1,
						EventTime:  e.EventTime,
					}).Error
				})
			}
		case domain.ProvideService:
			i := item.(*domain.ProvideServiceItem)
			if (i.From != nil && i.From.IsGroup()) || i.To.IsGroup() {
				g.Go(func() error {
					var activity string
					var serviceGiven int
					var serviceReceived int
					if i.From != nil && i.From.IsGroup() {
						serviceGiven = 1
						activity = fmt.Sprintf("group gave service '%s'", resourceMap[i.ResourceKey].Name)
					}
					if i.To.IsGroup() {
						serviceReceived = 1
						activity = fmt.Sprintf("group received service '%s'", resourceMap[i.ResourceKey].Name)
					}
					return h.db.Create(&readmodels.GroupReportItem{
						ID:               uuid.NewV4().String(),
						GroupKey:         e.GroupKey,
						Activity:         activity,
						GroupingID:       e.AggregateID,
						ServicesGiven:    serviceGiven,
						ServicesReceived: serviceReceived,
						EventTime:        e.EventTime,
					}).Error
				})
			}
		case domain.BorrowResource:
			i := item.(*domain.BorrowResourceItem)
			if i.To.IsGroup() {
				g.Go(func() error {
					return h.db.Create(&readmodels.GroupReportItem{
						ID:            uuid.NewV4().String(),
						GroupKey:      e.GroupKey,
						Activity:      fmt.Sprintf("group borrowed resource '%s'", resourceMap[i.ResourceKey].Name),
						GroupingID:    e.AggregateID,
						ItemsBorrowed: 1,
						EventTime:     e.EventTime,
					}).Error
				})
			}
			resource := resourceMap[i.ResourceKey]
			if resource.Owner.IsGroup() {
				g.Go(func() error {
					return h.db.Create(&readmodels.GroupReportItem{
						ID:         uuid.NewV4().String(),
						GroupKey:   resource.Owner.GetGroupKey(),
						Activity:   fmt.Sprintf("group lent resource '%s'", resourceMap[i.ResourceKey].Name),
						GroupingID: e.AggregateID,
						ItemsLent:  1,
						EventTime:  e.EventTime,
					}).Error
				})
			}
		}
	}

	return g.Wait()

}

func (h *GroupReportListener) handleResourceGroupSharingChanged(ctx context.Context, e resourcedomain.ResourceGroupSharingChanged) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, removedSharing := range e.ResourceGroupSharingChangedPayload.RemovedSharings {
		removedSharing := removedSharing
		g.Go(func() error {
			var offer int
			var request int
			if e.ResourceInfo.CallType == resourcedomain.Offer {
				offer = -1
			} else if e.ResourceInfo.CallType == resourcedomain.Request {
				request = -1
			}
			return h.db.Create(&readmodels.GroupReportItem{
				ID:            uuid.NewV4().String(),
				GroupKey:      removedSharing.GroupKey,
				Activity:      "Post removed from group",
				OfferCount:    offer,
				RequestsCount: request,
				EventTime:     e.EventTime,
			}).Error
		})
	}
	for _, addedSharing := range e.ResourceGroupSharingChangedPayload.AddedSharings {
		addedSharing := addedSharing
		g.Go(func() error {
			var offer int
			var request int
			if e.ResourceInfo.CallType == resourcedomain.Offer {
				offer = 1
			} else if e.ResourceInfo.CallType == resourcedomain.Request {
				request = 1
			}
			return h.db.Create(&readmodels.GroupReportItem{
				ID:            uuid.NewV4().String(),
				GroupKey:      addedSharing.GroupKey,
				Activity:      "Post added in group",
				OfferCount:    offer,
				RequestsCount: request,
				EventTime:     e.EventTime,
			}).Error
		})
	}
	return g.Wait()
}
