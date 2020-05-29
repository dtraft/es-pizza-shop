package eventsourcetest

import (
	"fmt"
	"reflect"
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"github.com/go-test/deep"
)

type HandleCommandTestCase struct {
	Label         string
	Given         []eventsource.EventData
	Command       eventsource.Command
	Expected      []eventsource.EventData
	ShouldError   bool
	ExpectedError error
}

type HandleCommandTestCases []*HandleCommandTestCase

func (c *HandleCommandTestCase) Test(a eventsource.Aggregate) error {
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
	if c.ShouldError {
		if c.ExpectedError != nil && c.ExpectedError != err {
			return fmt.Errorf("FAILED: %s.  Error: Expected error %#v, but got %#v", c.Label, c.ExpectedError, err)
		}

		if err == nil {
			return fmt.Errorf("FAILED: %s.  Error: Expected error, but got: %#v", c.Label, events)
		}
		return nil
	}

	for i, event := range events {
		exp := eventsource.NewEvent(a, c.Expected[i])
		got := eventsource.NewEvent(a, event)

		if exp.EventType != got.EventType {
			return fmt.Errorf("FAILED: %s at events[%d].  Error: Expected %s, got %s for EventType", c.Label, i, exp.EventType, got.EventType)
		}

		if exp.AggregateType != got.AggregateType {
			return fmt.Errorf("FAILED: %s at events[%d].  Expected %s, got %s for AggregateType", c.Label, i, exp.AggregateType, got.AggregateType)
		}

		if exp.EventTypeVersion != got.EventTypeVersion {
			return fmt.Errorf("FAILED: %s at events[%d].  Expected %d, got %d for EventTypeVersion", c.Label, i, exp.EventTypeVersion, got.EventTypeVersion)
		}

		if diff := deep.Equal(exp.Data, got.Data); diff != nil {
			return fmt.Errorf("FAILED: %s.  Error at events[%d]: %s", c.Label, i, diff)
		}
	}
	return nil
}

func (cases HandleCommandTestCases) Test(a eventsource.Aggregate, t *testing.T) {
	aggregateType, _ := eventsource.GetTypeName(a)
	for i, c := range cases {
		aggregate := reflect.New(aggregateType).Interface().(eventsource.Aggregate)
		if err := c.Test(aggregate); err != nil {
			t.Errorf("Case[%d] %s", i, err)
		}
	}
}

type ApplyEventTestCase struct {
	Given    []eventsource.EventData
	Event    eventsource.EventData
	Expected eventsource.Aggregate
}

func (c *ApplyEventTestCase) Test(a eventsource.Aggregate) error {
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
