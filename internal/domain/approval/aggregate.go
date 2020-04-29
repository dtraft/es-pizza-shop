package approval

import (
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
	case *RequestApproval:
		return a.handleRequestCommand(c)
	case *ReceiveApproval:
		return a.handleReceiveCommand(c)
	default:
		return nil, fmt.Errorf("No handler for command: %T", c)
	}
}

func (a *Aggregate) handleRequestCommand(c *RequestApproval) ([]eventsource.EventData, error) {
	if c.ApprovalID == 0 {
		return nil, fmt.Errorf("A valid approvalID was not provided, got: %d", c.ApprovalID)
	}

	if a.Sequence != 0 {
		return nil, nil
	}

	event := &ApprovalRequested{
		ApprovalID: c.ApprovalID,
	}
	return []eventsource.EventData{event}, nil
}

func (a *Aggregate) handleReceiveCommand(c *ReceiveApproval) ([]eventsource.EventData, error) {

	if a.Approved {
		return nil, nil
	}

	event := &ApprovalReceived{
		ApprovalID: c.ApprovalID,
	}
	return []eventsource.EventData{event}, nil
}

func (a *Aggregate) ApplyEvent(event eventsource.Event) error {

	switch event.Data.(type) {
	case *ApprovalReceived:
		a.Approved = true
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
