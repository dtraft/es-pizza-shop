package delivery

import (
	"errors"
	"fmt"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	. "forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/command"
	. "forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/event"
)

// Aggregate handles commands for orders
type Aggregate struct {
	eventsource.AggregateBase
	DeliveryID string
	Delivered  bool
}

func (a *Aggregate) Init(aggregateID string) {
	a.DeliveryID = aggregateID
}

// TestHandleCommand handles the commands for the Aggregate
func (a *Aggregate) HandleCommand(command eventsource.Command) ([]eventsource.EventData, error) {
	switch c := command.(type) {
	case *RequestDelivery:
		return a.handleRequestDelivery(c)
	case *ConfirmDelivery:
		return a.handleConfirmDelivery(c)
	default:
		message := fmt.Sprintf("No handler for command: %+v", c)
		return nil, errors.New(message)
	}
}
func (a *Aggregate) handleRequestDelivery(c *RequestDelivery) ([]eventsource.EventData, error) {
	if c.DeliveryID == 0 {
		return nil, fmt.Errorf("A valid approvalID was not provided, got: %d", c.DeliveryID)
	}

	if a.Sequence != 0 {
		return nil, nil
	}

	event := &DeliveryRequested{
		DeliveryID: c.DeliveryID,
	}
	return []eventsource.EventData{event}, nil
}

func (a *Aggregate) handleConfirmDelivery(c *ConfirmDelivery) ([]eventsource.EventData, error) {
	// If order has already been approved, don't emit any events
	if a.Delivered {
		return nil, nil
	}

	event := &DeliveryConfirmed{
		DeliveryID: c.DeliveryID,
	}
	return []eventsource.EventData{event}, nil
}

func (a *Aggregate) ApplyEvent(event eventsource.Event) error {

	switch event.Data.(type) {
	case *DeliveryConfirmed:
		a.Delivered = true
	}

	return nil
}

// AggregateID returns the AggregtateID
func (a *Aggregate) AggregateID() string {
	return a.DeliveryID
}

// ServiceType returns the AggregateType
func (a *Aggregate) Type() string {
	return "DeliveryAggregate"
}
