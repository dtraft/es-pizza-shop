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
		a.IncrementSequence()
	}

	// Handle command
	events, err := a.HandleCommand(c.Command)
	if err != nil {
		return err
	}

	for i, event := range events {
		exp := eventsource.NewEvent(a, c.Expected[i])
		got := eventsource.NewEvent(a, event)

		if exp.EventType != got.EventType {
			return fmt.Errorf("Expected %s, got %s for EventType", exp.EventType, got.EventType)
		}

		if exp.AggregateType != got.AggregateType {
			return fmt.Errorf("Expected %s, got %s for AggregateType", exp.AggregateType, got.AggregateType)
		}

		if exp.EventTypeVersion != got.EventTypeVersion {
			return fmt.Errorf("Expected %d, got %d for EventTypeVersion", exp.EventTypeVersion, got.EventTypeVersion)
		}

		if diff := deep.Equal(exp.Data, got.Data); diff != nil {
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
