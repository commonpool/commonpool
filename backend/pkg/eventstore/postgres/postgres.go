package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type PostgresEventStore struct {
	db *gorm.DB
}

func NewPostgresEventStore(db *gorm.DB) *PostgresEventStore {
	return &PostgresEventStore{
		db: db,
	}
}

var _ eventstore.EventStore = &PostgresEventStore{}

func (p PostgresEventStore) Load(ctx context.Context, streamKey eventstore.StreamKey) ([]*eventstore.StreamEvent, error) {
	var events []*eventstore.StreamEvent
	if err := p.db.
		Where("stream_type = ? AND stream_id = ?", streamKey.StreamType, streamKey.StreamID).
		Order("sequence_no asc").
		Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (p PostgresEventStore) Save(ctx context.Context, streamKey eventstore.StreamKey, expectedRevision int, events []*eventstore.StreamEvent) error {

	return p.db.Transaction(func(tx *gorm.DB) error {

		now := time.Now().UTC()

		correlationID := uuid.NewV4().String()
		correlationIDFromCtx := ctx.Value("correlationID")
		if correlationIDFromCtx != nil {
			if correlationIDStr, ok := correlationIDFromCtx.(string); ok {
				correlationID = correlationIDStr
			}
		}

		var stream eventstore.Stream
		query := tx.Model(eventstore.Stream{}).Find(&stream, "stream_id = ? and stream_type = ?", streamKey.StreamID, streamKey.StreamType)
		if err := query.Error; err != nil {
			return err
		}
		if query.RowsAffected == 0 {
			stream = eventstore.Stream{
				StreamID:      streamKey.StreamID,
				StreamType:    streamKey.StreamType,
				LatestVersion: 0,
			}
			if err := tx.Create(stream).Error; err != nil {
				return err
			}
		}

		if stream.LatestVersion != expectedRevision {
			return fmt.Errorf("could not save events: version mismatch: expected version %d but was %d", expectedRevision, stream.LatestVersion)
		}

		for i, event := range events {
			event.SequenceNo = expectedRevision + i
			if event.CorrelationID == "" {
				event.CorrelationID = correlationID
			}
			if (event.EventTime == time.Time{}) {
				event.EventTime = now
			} else {
				event.EventTime = event.EventTime.UTC()
			}
			if event.StreamKey() != streamKey {
				return fmt.Errorf("event streamKey != streamKey")
			}
			if event.EventID == "" {
				event.EventID = uuid.NewV4().String()
			}
		}

		if err := tx.Create(events).Error; err != nil {
			return fmt.Errorf("could not save events: %v", err)
		}

		stream.LatestVersion = stream.LatestVersion + len(events)

		if err := tx.Save(stream).Error; err != nil {
			return err
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (p PostgresEventStore) ReplayEventsByType(
	ctx context.Context,
	eventTypes []string,
	timestamp time.Time,
	replayFunc func(events []*eventstore.StreamEvent) error,
	options ...eventstore.ReplayEventsByTypeOptions) error {

	var streamEvents []*eventstore.StreamEvent

	batchSize := 200
	skip := 0

	if len(options) > 0 {
		option := options[0]
		if option.BatchSize > 0 {
			batchSize = option.BatchSize
		}
	}

	for {

		whereClause := "event_type in ("
		var params []interface{}
		for i, eventType := range eventTypes {
			whereClause += "?"
			if i < len(eventTypes)-1 {
				whereClause += ","
			}
			params = append(params, eventType)
		}
		whereClause += ") AND event_time > ?"
		params = append(params, timestamp.UTC())

		if err := p.db.
			Where(whereClause, params...).
			Order("event_time asc, sequence_no asc").
			Limit(batchSize).
			Offset(skip).
			Find(&streamEvents).Error; err != nil {
			return err
		}

		resultSize := len(streamEvents)

		if resultSize > 0 {
			if err := replayFunc(streamEvents); err != nil {
				return err
			}
		}

		if resultSize < batchSize {
			break
		} else {
			skip += batchSize
		}

	}

	return nil
}
