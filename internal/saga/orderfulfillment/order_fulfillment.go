package orderfulfillment

import (
	"encoding/json"
	"fmt"
	"strconv"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	orderEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
)

type OrderFulfillmentSaga struct {
	deliverySvc delivery.ServiceAPI
	approvalSvc approval.ServiceAPI

	Description string
	Approved    bool
	Delivered   bool
}

func (s *OrderFulfillmentSaga) Type() string {
	return "OrderFulfillmentSaga"
}

func (s *OrderFulfillmentSaga) Version() int {
	return 1
}

func (s *OrderFulfillmentSaga) StartEvent() string {
	return "OrderStartedEvent"
}

func (s *OrderFulfillmentSaga) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OrderFulfillmentSaga) AssociationID(event eventsource.Event) (*saga.SagaAssociation, error) {
	switch d := event.Data.(type) {
	case *orderEvents.OrderStartedEvent:
		return &saga.SagaAssociation{
			ID:              d.OrderID,
			AssociationType: "OrderID",
		}, nil
	case *orderEvents.OrderDescriptionSet:
		return &saga.SagaAssociation{
			ID:              d.OrderID,
			AssociationType: "OrderID",
		}, nil
	case *orderEvents.OrderSubmitted:
		return &saga.SagaAssociation{
			ID:              d.OrderID,
			AssociationType: "OrderID",
		}, nil
	default:
		return nil, fmt.Errorf("Unsupported event %T received: %+v", d, event)
	}
}

func (s *OrderFulfillmentSaga) HandleEvent(event eventsource.Event) (*saga.HandleEventResult, error) {
	switch d := event.Data.(type) {
	case *orderEvents.OrderStartedEvent:
		s.Description = d.Description
		return nil, nil
	case *orderEvents.OrderDescriptionSet:
		s.Description = d.Description
		return nil, nil
	case *orderEvents.OrderSubmitted:
		ids, err := s.startSaga(d, event)
		if err != nil {
			return nil, err
		}
		return &saga.HandleEventResult{AssociationIDs: ids}, nil
	default:
		return nil, fmt.Errorf("Unsupported event %T received: %+v", d, event)
	}
}

func (s *OrderFulfillmentSaga) startSaga(data *orderEvents.OrderSubmitted, event eventsource.Event) ([]*saga.SagaAssociation, error) {

	a, err := s.approvalSvc.SubmitOrderForApproval(&approval.OrderApproval{
		Description: s.Description,
	})
	if err != nil {
		return nil, err
	}
	associations := []*saga.SagaAssociation{
		{
			ID:              strconv.Itoa(a.ApprovalID),
			AssociationType: "ApprovalID",
		},
	}
	return associations, nil
}

var _ saga.SagaAPI = (*OrderFulfillmentSaga)(nil)
