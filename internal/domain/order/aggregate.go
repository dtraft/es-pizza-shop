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
	Status      Status
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
		return a.handleUpdateOrder(c)
	case *SubmitOrderCommand:
		return a.handleSubmitOrder(c)
	case *ApproveOrderCommand:
		return a.handleApproveOrder(c)
	case *DeliverOrderCommand:
		return a.handleDeliverOrder(c)
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

func (a *Aggregate) handleUpdateOrder(c *UpdateOrderCommand) ([]eventsource.EventData, error) {

	if a.Sequence == 0 {
		return nil, fmt.Errorf("No order found with id %s.", c.OrderID)
	}
	if a.Status != Started {
		return nil, fmt.Errorf("Cannot update an order which has already been submitted.")
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

func (a *Aggregate) handleSubmitOrder(c *SubmitOrderCommand) ([]eventsource.EventData, error) {

	if a.Sequence == 0 {
		return nil, fmt.Errorf("No order found with id %s.", c.OrderID)
	}

	if a.Status != Started {
		return nil, fmt.Errorf("Cannot submit order which has already been submitted.")
	}

	return []eventsource.EventData{
		&OrderSubmitted{OrderID: c.OrderID},
	}, nil
}

func (a *Aggregate) handleApproveOrder(c *ApproveOrderCommand) ([]eventsource.EventData, error) {

	if a.Sequence == 0 {
		return nil, fmt.Errorf("No order found with id %s.", c.OrderID)
	}

	if a.Status != Submitted {
		return nil, fmt.Errorf("Cannot approve order with status: %s.", a.Status)
	}

	return []eventsource.EventData{
		&OrderApproved{OrderID: c.OrderID},
	}, nil
}

func (a *Aggregate) handleDeliverOrder(c *DeliverOrderCommand) ([]eventsource.EventData, error) {

	if a.Sequence == 0 {
		return nil, fmt.Errorf("No order found with id %s.", c.OrderID)
	}

	if a.Status != Approved {
		return nil, fmt.Errorf("Cannot deliver order with status: %s.", a.Status)
	}

	return []eventsource.EventData{
		&OrderDelivered{OrderID: c.OrderID},
	}, nil
}

func (a *Aggregate) ApplyEvent(event eventsource.Event) error {

	switch e := event.Data.(type) {
	case *OrderStartedEvent:
		a.ServiceType = e.ServiceType
		a.Description = e.Description
		a.Status = Started
	case *OrderServiceTypeSetEvent:
		a.ServiceType = e.ServiceType
	case *OrderDescriptionSet:
		a.Description = e.Description
	case *OrderSubmitted:
		a.Status = Submitted
	case *OrderApproved:
		a.Status = Approved
	case *OrderDelivered:
		a.Status = Delivered
	default:
		return fmt.Errorf("Unsupported event %T received in ApplyEvent handler of the Order Aggregate: %+v", e, e)
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
