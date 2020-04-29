package approval

import (
	"errors"
	"fmt"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	. "forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/command"
	. "forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/event"
)

// Aggregate handles commands for orders
type Aggregate struct {
	eventsource.AggregateBase
	ApprovalID string
	Approved   bool
}

func (a *Aggregate) Init(aggregateID string) {
	a.ApprovalID = aggregateID
}

// TestHandleCommand handles the commands for the Aggregate
func (a *Aggregate) HandleCommand(command eventsource.Command) ([]eventsource.EventData, error) {
	switch c := command.(type) {
	case *ReceiveApproval:
		return a.handleReceiveCommand(c)
	default:
		message := fmt.Sprintf("No handler for command: %+v", c)
		return nil, errors.New(message)
	}
}

func (a *Aggregate) handleReceiveCommand(c *ReceiveApproval) ([]eventsource.EventData, error) {
	// If order has already been approved, don't emit any events
	if a.Sequence != 0 {
		return nil, nil
	}

	event := &ApprovalReceived{
		ApprovalID: c.ApprovalID,
	}
	return []eventsource.EventData{event}, nil
}

func (a *Aggregate) ApplyEvent(event eventsource.Event) error {

	switch e := event.Data.(type) {
	case *ApprovalReceived:
		a.Approved = true
	default:
		return fmt.Errorf("Unsupported event %T received in ApplyEvent handler of the Order Aggregate: %+v", e, e)
	}

	return nil
}

// AggregateID returns the AggregtateID
func (a *Aggregate) AggregateID() string {
	return a.ApprovalID
}

// ServiceType returns the AggregateType
func (a *Aggregate) Type() string {
	return "ApprovalAggregate"
}
