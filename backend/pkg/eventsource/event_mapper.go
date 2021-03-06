package eventsource

import "fmt"

type EventMapperFunc func(eventType string, bytes []byte) (Event, error)

type EventMapper struct {
	mappers map[string]EventMapperFunc
}

func NewEventMapper() *EventMapper {
	return &EventMapper{
		mappers: map[string]EventMapperFunc{},
	}
}

func (m *EventMapper) RegisterMapper(eventType string, mapper EventMapperFunc) error {
	if _, ok := m.mappers[eventType]; ok {
		return fmt.Errorf("mapper already registered for event type '%s'", eventType)
	}
	m.mappers[eventType] = mapper
	return nil
}

func (m *EventMapper) Map(eventType string, bytes []byte) (Event, error) {
	mapper, ok := m.mappers[eventType]
	if !ok {
		return nil, fmt.Errorf("mapper not registered for event type '%s'", eventType)
	}
	return mapper(eventType, bytes)
}
