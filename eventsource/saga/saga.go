package saga

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

type SagaManager struct {
	store Storer
}

func NewManager(store Storer) *SagaManager {
	return &SagaManager{
		store: store,
	}
}

func (m *SagaManager) ProcessEvent(event eventsource.Event, d SagaAPI) error {
	associationID, err := d.AssociationID(event)
	if err != nil {
		return err
	}

	// Load SagaWrapper
	w := &Wrapper{
		Type: d.Type(),
	}
	if d.StartEvent() == event.EventType {
		w.ID = uuid.New().String()
		if err := m.store.AddAssociationID(associationID, w); err != nil {
			return err
		}
	} else {
		w, err = m.store.Load(associationID, d.Type())
		if err != nil {
			return err
		}
		if err := d.Load(w.Data, w.Version); err != nil {
			return err
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

	w.Version = d.Version()
	w.Data = b

	if err := m.store.Save(w); err != nil {
		return err
	}

	if out != nil {
		for _, id := range out.AssociationIDs {
			if err := m.store.AddAssociationID(id, w); err != nil {
				return err
			}
		}
	}

	return nil
}

type Storer interface {
	Load(association *SagaAssociation, sagaType string) (*Wrapper, error)
	AddAssociationID(association *SagaAssociation, saga *Wrapper) error
	Save(saga *Wrapper) error
}

type SagaAPI interface {
	Type() string
	Version() int
	StartEvent() string
	Load(data json.RawMessage, version int) error

	AssociationID(event eventsource.Event) (*SagaAssociation, error)
	HandleEvent(event eventsource.Event) (*HandleEventResult, error)
}

type Wrapper struct {
	ID      string
	Version int
	Type    string
	Data    json.RawMessage
}

type HandleEventResult struct {
	AssociationIDs []*SagaAssociation
}

type SagaAssociation struct {
	ID              string
	AssociationType string
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
