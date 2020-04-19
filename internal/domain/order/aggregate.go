package order

import (
	"errors"
	"fmt"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	. "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/command"
	. "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
	. "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
)

// Aggregate handles commands for orders
type Aggregate struct {
	OrderID     string
	ServiceType ServiceType
}

func (a *Aggregate) Init(aggregateID string) {
	a.OrderID = aggregateID
}

// HandleCommand handles the commands for the Aggregate
func (a *Aggregate) HandleCommand(command eventsource.Command) ([]eventsource.Event, error) {
	switch c := command.(type) {
	case *StartOrderCommand:
		event := &OrderStartedEvent{
			OrderID:     c.OrderID,
			ServiceType: c.Type,
		}
		return []eventsource.Event{eventsource.NewEvent(a, event)}, nil
	case *ToggleOrderServiceTypeCommand:

		var serviceType ServiceType
		if a.ServiceType == Pickup {
			serviceType = Delivery
		} else {
			serviceType = Pickup
		}

		event := &OrderServiceTypeSetEvent{
			OrderID:     c.OrderID,
			ServiceType: serviceType,
		}
		return []eventsource.Event{eventsource.NewEvent(a, event)}, nil
	default:
		message := fmt.Sprintf("No handler for command type: %+v", c)
		return nil, errors.New(message)
	}

}

func (a *Aggregate) ApplyEvent(event eventsource.Event) error {

	switch e := event.Data.(type) {
	case *OrderStartedEvent:
		a.ServiceType = e.ServiceType
	case *OrderServiceTypeSetEvent:
		a.ServiceType = e.ServiceType
	default:
		return fmt.Errorf("Unsupported event received in ApplyEvent handler of the Order Aggregate: %+v", e)
	}

	return nil
}

// AggregateID returns the AggregtateID
func (a *Aggregate) AggregateID() string {
	return a.OrderID
}

// Type returns the AggregateType
func (a *Aggregate) Type() string {
	return "OrderAggregate"
}
