package eventsourcetest

import (
	"fmt"
	"reflect"
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"

	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

type SagaLoadTestCase struct {
	Label         string
	Version       int
	Saga          string
	Expected      saga.SagaAPI
	ShouldError   bool
	ExpectedError error
}

type SagaLoadTestCases []*SagaLoadTestCase

func (c *SagaLoadTestCase) Test() error {
	sagaType, _ := eventsource.GetTypeName(c.Expected)
	got := reflect.New(sagaType).Interface().(saga.SagaAPI)
	err := got.Load([]byte(c.Saga), c.Version)

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

func (cases SagaLoadTestCases) Test(t *testing.T) {
	for i, c := range cases {
		if err := c.Test(); err != nil {
			t.Errorf("Case[%d] %s", i, err)
		}
	}
}

type SagaAssociationIDTestCase struct {
	Label         string
	Saga          saga.SagaAPI
	Event         eventsource.Event
	Expected      *saga.SagaAssociation
	ShouldError   bool
	ExpectedError error
}

type SagaAssociationIDTestCases []*SagaAssociationIDTestCase

func (c *SagaAssociationIDTestCase) Test() error {
	association, err := c.Saga.AssociationID(c.Event)

	if c.ShouldError {
		if diff := deep.Equal(c.ExpectedError, err); c.ExpectedError != nil && diff != nil {
			return fmt.Errorf("FAILED: %s.  Error: error does not match expected, details: %s", c.Label, diff)
		}

		if err == nil {
			return fmt.Errorf("FAILED: %s.  Error: Expected error, but got: %#v", c.Label, association)
		}
	} else {
		if diff := deep.Equal(association, c.Expected); diff != nil {
			return fmt.Errorf("FAILED: %s.  Error: %s", c.Label, diff)
		}
	}

	return nil
}

func (cases SagaAssociationIDTestCases) Test(t *testing.T) {
	for i, c := range cases {
		if err := c.Test(); err != nil {
			t.Errorf("Case[%d] %s", i, err)
		}
	}
}

type SagaHandleEventTestCase struct {
	Label          string
	Given          []eventsource.Event
	Event          eventsource.Event
	Saga           saga.SagaAPI
	ExpectedSaga   saga.SagaAPI
	ExpectedResult *saga.HandleEventResult
	ShouldError    bool
	ExpectedError  error
}

type SagaHandleEventTestCases []*SagaHandleEventTestCase

func (c *SagaHandleEventTestCase) Test() error {
	var got saga.SagaAPI
	if c.Saga != nil {
		got = c.Saga
	} else {
		sagaType, _ := eventsource.GetTypeName(c.ExpectedSaga)
		got = reflect.New(sagaType).Interface().(saga.SagaAPI)
	}

	for _, e := range c.Given {
		got.HandleEvent(e)
	}
	result, err := got.HandleEvent(c.Event)

	if c.ShouldError {
		if diff := deep.Equal(c.ExpectedError, err); c.ExpectedError != nil && diff != nil {
			return fmt.Errorf("FAILED: %s.  Error: error does not match expected, details: %s", c.Label, diff)
		}

		if err == nil {
			return fmt.Errorf("FAILED: %s.  Error: Expected error, but got: %#v", c.Label, result)
		}
	} else {
		if diff := deep.Equal(result, c.ExpectedResult); diff != nil {
			return fmt.Errorf("FAILED: %s.  Error: result does not match expected, details: %s", c.Label, diff)
		}
		if diff := deep.Equal(got, c.ExpectedSaga); diff != nil {
			return fmt.Errorf("FAILED: %s.  Error: saga does not match expected, details: %s", c.Label, diff)
		}
	}

	return nil
}

func (cases SagaHandleEventTestCases) Test(t *testing.T) {
	for i, c := range cases {
		if err := c.Test(); err != nil {
			t.Errorf("Case[%d] %s", i, err)
		}
	}
}
