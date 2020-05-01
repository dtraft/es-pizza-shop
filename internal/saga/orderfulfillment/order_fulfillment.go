package orderfulfillment

import (
	"fmt"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"
	orderEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
)

type OrderFulfillmentSaga struct {
	saga.BaseSaga
	sagaStore    saga.SagaStorer
	orderService order.ServiceAPI

	Approved  bool
	Delivered bool
}

func (s *OrderFulfillmentSaga) HandleEvent(event eventsource.Event) error {
	switch d := event.Data.(type) {
	case *orderEvents.OrderStartedEvent:
		return s.handleOrderStartedEvent(d, event)
	default:
		return fmt.Errorf("Unsupported event %T received: %+v", d, event)
	}
}

func (s *OrderFulfillmentSaga) handleOrderStartedEvent(data *orderEvents.OrderStartedEvent, event eventsource.Event) error {
	return nil
}
