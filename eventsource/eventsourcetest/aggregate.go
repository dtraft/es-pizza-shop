package eventsourcetest

import (
	"fmt"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"github.com/go-test/deep"
)

type HandleCommandCase struct {
	Given    []eventsource.EventData
	Command  eventsource.Command
	Expected []eventsource.EventData
}

type ApplyEventCase struct {
	Given    []eventsource.EventData
	Event    eventsource.EventData
	Expected eventsource.Aggregate
}

func (c *HandleCommandCase) TestHandleCommand(a eventsource.Aggregate) error {
	// Setup aggregate for testing
	for _, e := range c.Given {
		event := eventsource.NewEvent(a, e)
		if err := a.ApplyEvent(event); err != nil {
			return err
		}
	}

	// Handle command
	events, err := a.HandleCommand(c.Command)
	if err != nil {
		return err
	}

	for i, event := range events {
		exp := eventsource.NewEvent(a, c.Expected[i])

		if exp.EventType != event.EventType {
			return fmt.Errorf("Expected %s, got %s for EventType", exp.EventType, event.EventType)
		}

		if exp.AggregateType != event.AggregateType {
			return fmt.Errorf("Expected %s, got %s for AggregateType", exp.AggregateType, event.AggregateType)
		}

		if exp.EventTypeVersion != event.EventTypeVersion {
			return fmt.Errorf("Expected %d, got %d for EventTypeVersion", exp.EventTypeVersion, event.EventTypeVersion)
		}

		if diff := deep.Equal(exp.Data, event.Data); diff != nil {
			return fmt.Errorf("Events[%d]: %s", i, diff)
		}
	}
	return nil
}

func (c *ApplyEventCase) TestApplyEvent(a eventsource.Aggregate) error {
	// Setup aggregate for testing
	for _, e := range c.Given {
		event := eventsource.NewEvent(a, e)
		if err := a.ApplyEvent(event); err != nil {
			return err
		}
	}

	event := eventsource.NewEvent(a, c.Event)
	if err := a.ApplyEvent(event); err != nil {
		return err
	}

	if diff := deep.Equal(a, c.Expected); diff != nil {
		return fmt.Errorf("%s", diff)
	}

	return nil
}
