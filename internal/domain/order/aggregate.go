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

// HandleCommand handles the commands for the Aggregate
func (a *Aggregate) HandleCommand(command eventsource.Command) ([]eventsource.Event, error) {
	switch c := command.(type) {
	case *StartOrderCommand:
		return a.handleStartOrder(c)
	case *ToggleOrderServiceTypeCommand:
		return a.handleToggleServiceType(c)
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
		ServiceType: c.Type,
		Description: c.Description,
	}
	return []eventsource.Event{eventsource.NewEvent(a, event)}, nil
}

func (a *Aggregate) handleToggleServiceType(c *ToggleOrderServiceTypeCommand) ([]eventsource.Event, error) {
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
}

func (a *Aggregate) handleOrderUpdate(c *UpdateOrderCommand) ([]eventsource.Event, error) {
	var events []eventsource.Event

	// Service Type
	if c.ServiceType != a.ServiceType {
		event := &OrderServiceTypeSetEvent{
			OrderID:     c.OrderID,
			ServiceType: c.ServiceType,
		}
		events = append(events, eventsource.NewEvent(a, event))
	}

	// Description
	if c.Description != a.Description {
		event := &OrderDescriptionSet{
			OrderID:     c.OrderID,
			Description: c.Description,
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
