package eventsource

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Command interface for commands in the domain model
type Command interface {
	AggregateID() string
}

// Event stores the data for every event
type Event struct {
	EventID          string      `json:"eventId"`
	AggregateID      string      `json:"aggregateId"`
	AggregateType    string      `json:"aggregateType"`
	EventTypeVersion int         `json:"eventVersion"`
	EventType        string      `json:"eventType"`
	Timestamp        time.Time   `json:"eventTimestamp"`
	Data             interface{} `json:"eventData"`
}

// EventData interface for commands in the domain model
type EventData interface {
	Load(json.RawMessage, int) error
	Version() int
}

// Aggregate interace for Aggregates in the domain model
type Aggregate interface {
	Init(aggregateID string)
	AggregateID() string
	Type() string
	HandleCommand(command Command) ([]Event, error)
	ApplyEvent(event Event) error
}

type Projection interface {
	ApplyEvent(event Event) error
}

type EventStore interface {
	SaveEvent(event Event) error
	EventsForAggregate(aggregateID string) ([]Event, error)
}

type EventSource struct {
	store EventStore
}

func New(eventStore EventStore) EventSource {
	return EventSource{
		store: eventStore,
	}
}

func (es *EventSource) LoadAggregate(a Aggregate) error {
	events, err := es.store.EventsForAggregate(a.AggregateID())
	if err != nil {
		return err
	}

	for _, event := range events {
		if err = a.ApplyEvent(event); err != nil {
			return err
		}
	}

	return nil
}

func (es *EventSource) ProcessCommand(c Command, a Aggregate) error {
	// Restore the aggregate
	a.Init(c.AggregateID())
	if err := es.LoadAggregate(a); err != nil {
		return err
	}

	events, err := a.HandleCommand(c)
	if err != nil {
		return err
	}
	for _, event := range events {
		err := a.ApplyEvent(event)
		if err != nil {
			return err
		}

		err = es.store.SaveEvent(event)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewEvent publishes the event
func NewEvent(a Aggregate, p EventData) Event {
	_, eventType := GetTypeName(p)
	return Event{
		EventID:          uuid.New().String(),
		AggregateID:      a.AggregateID(),
		AggregateType:    a.Type(),
		EventTypeVersion: p.Version(),
		EventType:        eventType,
		Timestamp:        time.Now(),
		Data:             p,
	}
}

func (e *Event) Load(data []byte) error {
	eventType, err := GetEventOfType(e.EventType)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, e); err != nil {
		return err
	}

	r, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}

	if err := eventType.Load(r, e.EventTypeVersion); err != nil {
		return err
	}

	e.Data = eventType

	return nil
}
