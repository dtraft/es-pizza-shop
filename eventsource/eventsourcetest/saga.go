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
	eventType, _ := eventsource.GetTypeName(c.Expected)
	got := reflect.New(eventType).Interface().(eventsource.EventData)
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
	Expected      saga.SagaAssociation
	ShouldError   bool
	ExpectedError error
}

type SagaAssociationIDTestCases []*SagaAssociationIDTestCase

func (c *SagaAssociationIDTestCase) Test() error {
	association, err := c.Saga.AssociationID(c.Event)

	if c.ShouldError {
		if c.ExpectedError != nil && c.ExpectedError != err {
			return fmt.Errorf("FAILED: %s.  Error: Expected error %#v, but got %#v", c.Label, c.ExpectedError, err)
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
