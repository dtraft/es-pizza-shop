package orderfulfillment

import (
	"fmt"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"
	orderEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
)

type OrderFulfillmentSaga struct {
	saga.BaseSaga
	deliverySvc delivery.ServiceAPI
	approvalSvc approval.ServiceAPI

	Description string
	Approved    bool
	Delivered   bool
}

func (s *OrderFulfillmentSaga) SagaType() string {
	return "OrderFulfillmentSaga"
}

func (s *OrderFulfillmentSaga) StartSagaEventTypes() []string {
	return []string{"OrderStartedEvent"}
}

func (s *OrderFulfillmentSaga) HandleEvent(event eventsource.Event) error {
	switch d := event.Data.(type) {
	case *orderEvents.OrderStartedEvent:
		s.Description = d.Description
		return nil
	case *orderEvents.OrderDescriptionSet:
		s.Description = d.Description
		return nil
	case *orderEvents.OrderSubmitted:
		return s.startSaga(d, event)
	default:
		return fmt.Errorf("Unsupported event %T received: %+v", d, event)
	}
}

func (s *OrderFulfillmentSaga) startSaga(data *orderEvents.OrderSubmitted, event eventsource.Event) error {

	s.approvalSvc.SubmitOrderForApproval(&approval.OrderApproval{
		Description: "",
	})
	return nil
}
