package saga

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

type SagaManager struct {
	store SagaStorer
}

func (m *SagaManager) ProcessEvent(event eventsource.Event, d SagaAPI) error {
	associationID, err := d.AssociationID(event)
	if err != nil {
		return err
	}

	// Load SagaWrapper
	var s *SagaWrapper
	if d.StartEvent() == event.EventType {
		s = &SagaWrapper{
			ID:   uuid.New().String(),
			Data: d,
		}
		if err := m.store.AddAssociationID(associationID, s); err != nil {
			return err
		}

	} else {
		raw, err := m.store.Load(associationID, d.Type())
		if err != nil {
			return err
		}

		if err := d.Load(raw.Data, raw.Version); err != nil {
			return err
		}

		s = &SagaWrapper{
			ID:   raw.ID,
			Data: d,
		}
	}

	// Handle Event
	out, err := d.HandleEvent(event)
	if err != nil {
		return err
	}

	// Save SagaWrapper
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	raw := &RawSagaWrapper{
		ID:      s.ID,
		Version: d.Version(),
		Type:    d.Type(),
		Data:    b,
	}
	if err := m.store.Save(raw); err != nil {
		return err
	}

	if out != nil {
		for _, id := range out.AssociationIDs {
			if err := m.store.AddAssociationID(id, raw); err != nil {
				return err
			}
		}
	}

	return nil
}

type SagaStorer interface {
	Load(associationID string, sagaType string) (*RawSagaWrapper, error)
	AddAssociationID(associationID string, saga *RawSagaWrapper) error
	Save(saga *RawSagaWrapper) error
}

type SagaAPI interface {
	Type() string
	Version() int
	StartEvent() string
	Load(data json.RawMessage, version int) error

	AssociationID(event eventsource.Event) (string, error)
	HandleEvent(event eventsource.Event) (*HandleEventResult, error)
}
type SagaWrapper struct {
	ID   string
	Data SagaAPI
}

type RawSagaWrapper struct {
	ID      string
	Version int
	Type    string
	Data    json.RawMessage
}

type HandleEventResult struct {
	AssociationIDs []string
}

type SagaAssociationNotFoundError struct {
	AssociationID string
	SagaType      string
}

func (e *SagaAssociationNotFoundError) Error() string {
	return fmt.Sprintf("No %s saga found for AssociationID %s", e.SagaType, e.AssociationID)
}

type SagaNotFoundError struct {
	SagaID string
}

func (e *SagaNotFoundError) Error() string {
	return fmt.Sprintf("No saga found for SagaID %s", e.SagaID)
}
