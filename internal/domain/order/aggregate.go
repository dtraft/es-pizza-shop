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
	Description string
}

func (a *Aggregate) Init(aggregateID string) {
	a.OrderID = aggregateID
}

// TestHandleCommand handles the commands for the Aggregate
func (a *Aggregate) HandleCommand(command eventsource.Command) ([]eventsource.Event, error) {
	switch c := command.(type) {
	case *StartOrderCommand:
		return a.handleStartOrder(c)
	case *UpdateOrderCommand:
		return a.handleOrderUpdate(c)
	default:
		message := fmt.Sprintf("No handler for command: %+v", c)
		return nil, errors.New(message)
	}
}

func (a *Aggregate) handleStartOrder(c *StartOrderCommand) ([]eventsource.Event, error) {
	event := &OrderStartedEvent{
		OrderID:     c.OrderID,
		ServiceType: c.ServiceType,
		Description: c.Description,
	}
	return []eventsource.Event{eventsource.NewEvent(a, event)}, nil
}

func (a *Aggregate) handleOrderUpdate(c *UpdateOrderCommand) ([]eventsource.Event, error) {
	var events []eventsource.Event

	// Service ServiceType
	serviceType, err := c.ServiceType.Get()
	if err == nil && serviceType != a.ServiceType {
		event := &OrderServiceTypeSetEvent{
			OrderID:     c.OrderID,
			ServiceType: serviceType,
		}
		events = append(events, eventsource.NewEvent(a, event))
	}

	// Description
	description, err := c.Description.Get()
	if err == nil && description != a.Description {
		event := &OrderDescriptionSet{
			OrderID:     c.OrderID,
			Description: description,
		}
		events = append(events, eventsource.NewEvent(a, event))
	}

	return events, nil
}

func (a *Aggregate) ApplyEvent(event eventsource.Event) error {

	switch e := event.Data.(type) {
	case *OrderStartedEvent:
		a.ServiceType = e.ServiceType
		a.Description = e.Description
	case *OrderServiceTypeSetEvent:
		a.ServiceType = e.ServiceType
	case *OrderDescriptionSet:
		a.Description = e.Description
	default:
		return fmt.Errorf("Unsupported event received in TestApplyEvent handler of the Order Aggregate: %+v", e)
	}

	return nil
}

// AggregateID returns the AggregtateID
func (a *Aggregate) AggregateID() string {
	return a.OrderID
}

// ServiceType returns the AggregateType
func (a *Aggregate) Type() string {
	return "OrderAggregate"
}
