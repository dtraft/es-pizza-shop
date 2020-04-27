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
	eventsource.AggregateBase
	OrderID     string
	ServiceType ServiceType
	Description string
}

func (a *Aggregate) Init(aggregateID string) {
	a.OrderID = aggregateID
}

// TestHandleCommand handles the commands for the Aggregate
func (a *Aggregate) HandleCommand(command eventsource.Command) ([]eventsource.EventData, error) {
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

func (a *Aggregate) handleStartOrder(c *StartOrderCommand) ([]eventsource.EventData, error) {

	if a.Sequence != 0 {
		return nil, fmt.Errorf("An order with id %s already exists.", c.OrderID)
	}

	event := &OrderStartedEvent{
		OrderID:     c.OrderID,
		ServiceType: c.ServiceType,
		Description: c.Description,
	}
	return []eventsource.EventData{event}, nil
}

func (a *Aggregate) handleOrderUpdate(c *UpdateOrderCommand) ([]eventsource.EventData, error) {

	if a.Sequence == 0 {
		return nil, fmt.Errorf("No order found with id %s.", c.OrderID)
	}

	var events []eventsource.EventData

	// Service ServiceType
	serviceType, err := c.ServiceType.Get()
	if err == nil && serviceType != a.ServiceType {
		event := &OrderServiceTypeSetEvent{
			OrderID:     c.OrderID,
			ServiceType: serviceType,
		}
		events = append(events, event)
	}

	// Description
	description, err := c.Description.Get()
	if err == nil && description != a.Description {
		event := &OrderDescriptionSet{
			OrderID:     c.OrderID,
			Description: description,
		}
		events = append(events, event)
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
