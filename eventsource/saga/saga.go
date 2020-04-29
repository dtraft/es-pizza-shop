package saga

import "fmt"

type SagaStorer interface {
	Load(associationID string, sagaType string, in interface{}) error
	AddAssociationID(associationID string, sagaType string) error
	RemoveAssociationID(associationID string, sagaType string) error
	Save(saga SagaAPI) error
}

type SagaAPI interface {
	SagaID() string
}

type BaseSaga struct {
	SagaID string
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
