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
	EventID           string      `json:"eventId"`
	AggregateID       string      `json:"aggregateId"`
	AggregateType     string      `json:"aggregateType"`
	AggregateSequence int         `json:"aggregateSequence"`
	EventTypeVersion  int         `json:"eventVersion"`
	EventType         string      `json:"eventType"`
	Timestamp         time.Time   `json:"eventTimestamp"`
	Data              interface{} `json:"eventData"`
}

type eventIntermediate struct {
	*Event
	Data json.RawMessage `json:"eventData"`
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
	HandleCommand(command Command) ([]EventData, error)
	ApplyEvent(event Event) error
	setSequence(int)
	getSequence() int
	IncrementSequence()
}

type AggregateBase struct {
	Sequence int
}

// IncrementVersion ads 1 to the current version
func (b *AggregateBase) IncrementSequence() {
	b.Sequence++
}

func (b *AggregateBase) setSequence(seq int) {
	b.Sequence = seq
}

func (b *AggregateBase) getSequence() int {
	return b.Sequence
}

type Projection interface {
	HandleEvent(event Event) error
}

type EventStorer interface {
	SaveEvent(event Event) error
	EventsForAggregate(aggregateID string) ([]Event, error)
}

type EventSourceAPI interface {
	LoadAggregate(a Aggregate) error
	ProcessCommand(c Command, a Aggregate) error
}

type EventSource struct {
	store EventStorer
}

func New(eventStore EventStorer) *EventSource {
	return &EventSource{
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
		a.setSequence(event.AggregateSequence)
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
		a.IncrementSequence()
		e := NewEvent(a, event)
		err = es.store.SaveEvent(e)
		// TODO - implement retries if we get an AggregateLockError here
		if err != nil {
			return err
		}

		err := a.ApplyEvent(e)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewEvent publishes the event
func NewEvent(a Aggregate, p EventData) Event {
	_, eventType := GetTypeName(p)
	event := Event{
		EventID:           uuid.New().String(),
		AggregateID:       a.AggregateID(),
		AggregateType:     a.Type(),
		AggregateSequence: a.getSequence(),
		EventType:         eventType,
		EventTypeVersion:  p.Version(),
		Timestamp:         time.Now(),
		Data:              p,
	}
	return event
}

func (e *Event) Load(data []byte) error {
	eventType, err := GetEventOfType(e.EventType)
	if err != nil {
		return err
	}

	temp := &eventIntermediate{
		Event: e,
	}
	if err = json.Unmarshal(data, temp); err != nil {
		return err
	}

	if err := eventType.Load(temp.Data, e.EventTypeVersion); err != nil {
		return err
	}

	e.Data = eventType

	return nil
}
