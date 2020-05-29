package eventsourcetest

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

type EventLoadTestCase struct {
	Label         string
	Version       int
	Event         string
	Expected      eventsource.EventData
	ShouldError   bool
	ExpectedError error
}

type EventLoadTestCases []*EventLoadTestCase

func (c *EventLoadTestCase) Test() error {
	eventType, _ := eventsource.GetTypeName(c.Expected)
	got := reflect.New(eventType).Interface().(eventsource.EventData)
	err := got.Load([]byte(c.Event), c.Version)

	if c.ShouldError {
		if c.ExpectedError != nil && c.ExpectedError != err {
			return fmt.Errorf("FAILED: %s.  Error: Expected error %#v, but got %#v", c.Label, c.ExpectedError, err)
		}

		if err == nil {
			return fmt.Errorf("FAILED: %s.  Error: Expected error, but got: %#v", c.Label, got)
		}
	} else {
		if diff := deep.Equal(got, c.Expected); diff != nil {
			return fmt.Errorf("FAILED: %s.  Error: %s", c.Label, diff)
		}
	}

	return nil
}

func (cases EventLoadTestCases) Test(t *testing.T) {
	for i, c := range cases {
		if err := c.Test(); err != nil {
			t.Errorf("Case[%d] %s", i, err)
		}
	}
}
