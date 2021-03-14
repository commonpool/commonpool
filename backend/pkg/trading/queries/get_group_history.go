package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	groupreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"gorm.io/gorm"
)

type GetGroupHistory struct {
	db *gorm.DB
}

func NewGetGroupHistory(db *gorm.DB) *GetGroupHistory {
	return &GetGroupHistory{db: db}
}

func (q *GetGroupHistory) Get(ctx context.Context, groupKey keys.GroupKey) ([]*groupreadmodels.GroupReportItem, error) {

	var entries []*groupreadmodels.GroupReportItem
	if err := q.db.Model(&groupreadmodels.GroupReportItem{}).
		Where("group_key = ?", groupKey).
		Order("event_time asc").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	var cumulativeItems = []*groupreadmodels.GroupReportItem{}
	for i, entry := range entries {
		if i == 0 {
			cumulativeItems = append(cumulativeItems, entry)
			continue
		}
		lastEntry := cumulativeItems[i-1]

		newEntry := &groupreadmodels.GroupReportItem{
			ID:               entry.ID,
			GroupKey:         entry.GroupKey,
			Activity:         entry.Activity,
			GroupingID:       entry.GroupingID,
			ItemsReceived:    lastEntry.ItemsReceived + entry.ItemsReceived,
			ItemsGiven:       lastEntry.ItemsGiven + entry.ItemsGiven,
			ItemsOwned:       lastEntry.ItemsOwned + entry.ItemsOwned,
			ItemsLent:        lastEntry.ItemsLent + entry.ItemsLent,
			ItemsBorrowed:    lastEntry.ItemsBorrowed + entry.ItemsBorrowed,
			ServicesGiven:    lastEntry.ServicesGiven + entry.ServicesGiven,
			ServicesReceived: lastEntry.ServicesReceived + entry.ServicesReceived,
			OfferCount:       lastEntry.OfferCount + entry.OfferCount,
			RequestsCount:    lastEntry.RequestsCount + entry.RequestsCount,
			HoursInBank:      lastEntry.HoursInBank + entry.HoursInBank,
			EventTime:        entry.EventTime,
		}
		cumulativeItems = append(cumulativeItems, newEntry)
	}

	return cumulativeItems, nil

}
