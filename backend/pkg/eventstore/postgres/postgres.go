package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type PostgresEventStore struct {
	db          *gorm.DB
	eventMapper *eventsource.EventMapper
}

func NewPostgresEventStore(db *gorm.DB, eventMapper *eventsource.EventMapper) *PostgresEventStore {
	return &PostgresEventStore{
		db:          db,
		eventMapper: eventMapper,
	}
}

var _ eventstore.EventStore = &PostgresEventStore{}

func (p *PostgresEventStore) MigrateDatabase() error {
	return p.db.AutoMigrate(&eventstore.StreamEvent{}, &eventstore.Stream{})
}

func (p PostgresEventStore) Load(ctx context.Context, streamKey eventstore.StreamKey) ([]eventsource.Event, error) {
	var events []*eventstore.StreamEvent
	if err := p.db.
		Where("stream_type = ? AND stream_id = ?", streamKey.StreamType, streamKey.StreamID).
		Order("sequence_no asc").
		Find(&events).Error; err != nil {
		return nil, err
	}

	mappedEvents, err := p.mapEvents(events)
	if err != nil {
		return nil, err
	}

	return mappedEvents, nil
}

func (p PostgresEventStore) mapEvents(events []*eventstore.StreamEvent) ([]eventsource.Event, error) {
	var mappedEvents = make([]eventsource.Event, len(events))
	for i, event := range events {
		mappedEvent, err := p.eventMapper.Map(event.EventType, []byte(event.Body))
		if err != nil {
			return nil, err
		}
		mappedEvents[i] = mappedEvent
	}
	return mappedEvents, nil
}

func eventStringValueOrDefault(defaultValue string, key string, tempStruct map[string]interface{}) string {
	var value = defaultValue
	tempValue, ok := tempStruct[key]
	if !ok {
		tempStruct[key] = value
	} else {
		valueStr, ok := tempValue.(string)
		if !ok {
			tempStruct[key] = value
		} else if valueStr == "" {
			tempStruct[key] = value
		} else {
			value = valueStr
		}
	}
	return value
}

func eventTimeValueOrDefault(defaultValue time.Time, key string, tempStruct map[string]interface{}) time.Time {
	var value = defaultValue
	tempValue, ok := tempStruct[key]
	if !ok {
		tempStruct[key] = value
	} else {
		valueStr, ok := tempValue.(string)
		if !ok {
			tempStruct[key] = value
			return value
		}

		valueTime, err := time.Parse(time.RFC3339Nano, valueStr)
		if err != nil {
			tempStruct[key] = value
			return value
		}

		if (valueTime == time.Time{}) {
			tempStruct[key] = value
			return value
		}

		if valueTime != valueTime.UTC() {
			tempStruct[key] = valueTime.UTC()
			return valueTime.UTC()
		}

	}
	return value
}

func (p PostgresEventStore) Save(ctx context.Context, streamKey eventstore.StreamKey, expectedRevision int, events []eventsource.Event) error {

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

		var streamEvents = make([]*eventstore.StreamEvent, len(events))

		for i, event := range events {

			evtJson, err := json.Marshal(event)
			if err != nil {
				return err
			}

			var tempStruct map[string]interface{}
			err = json.Unmarshal(evtJson, &tempStruct)
			if err != nil {
				return err
			}

			evtCorrelationId := eventStringValueOrDefault(correlationID, "correlation_id", tempStruct)
			evtAggregateType := eventStringValueOrDefault(streamKey.StreamType, "aggregate_type", tempStruct)
			evtAggregateId := eventStringValueOrDefault(streamKey.StreamID, "aggregate_id", tempStruct)
			evtEventId := eventStringValueOrDefault(uuid.NewV4().String(), "event_id", tempStruct)
			evtTime := eventTimeValueOrDefault(now, "event_time", tempStruct)
			evtRevision := expectedRevision + i
			tempStruct["sequence_no"] = evtRevision

			eventBody, err := json.Marshal(tempStruct)
			if err != nil {
				return err
			}

			streamEvent := &eventstore.StreamEvent{
				SequenceNo:    evtRevision,
				EventTime:     evtTime,
				CorrelationID: evtCorrelationId,
				StreamID:      evtAggregateId,
				StreamType:    evtAggregateType,
				EventID:       evtEventId,
				EventType:     event.GetEventType(),
				EventVersion:  event.GetEventVersion(),
				Body:          string(eventBody),
			}

			if streamEvent.StreamKey() != streamKey {
				return fmt.Errorf("event streamKey != streamKey")
			}

			streamEvents[i] = streamEvent
		}

		if err := tx.Create(streamEvents).Error; err != nil {
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
	replayFunc func(events []eventsource.Event) error,
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

			evts, err := p.mapEvents(streamEvents)
			if err != nil {
				return err
			}
			if err := replayFunc(evts); err != nil {
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
