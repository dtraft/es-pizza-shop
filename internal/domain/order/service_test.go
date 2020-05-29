package order

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/markphelps/optional"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/command"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

func TestService_StartOrder(t *testing.T) {
	cases := []struct {
		Label       string
		Order       *model.Order
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should assign a UUID if orderId is not provided",
			Order: &model.Order{},
			Check: func(c eventsource.Command) error {
				if c.AggregateID() == "" {
					return fmt.Errorf("OrderID is empty.")
				}
				return nil
			},
		},
		{
			Label: "Should correctly issue the start order command",
			Order: &model.Order{
				OrderID:     "testOrderId",
				ServiceType: model.Pickup,
				Description: "I'm a test!",
			},
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.StartOrderCommand)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.StartOrderCommand{}, c)
				}
				if cmd.OrderID != "testOrderId" {
					return fmt.Errorf("Expected `%s` for OrderID, got `%s`", "testOrderId", cmd.OrderID)
				}
				if cmd.ServiceType != model.Pickup {
					return fmt.Errorf("Expected `%d` for ServiceType, got `%d`", model.Pickup, cmd.ServiceType)
				}
				if cmd.Description != "I'm a test!" {
					return fmt.Errorf("Expected `%s` for Description, got `%s`", "I'm a test!", cmd.Description)
				}
				return nil
			},
		},
		{
			Label: "Should bubble up errors",
			Order: &model.Order{},
			Check: func(c eventsource.Command) error {
				return nil
			},
			ShouldError: true,
		},
	}

	for i, c := range cases {
		s := NewService(&mockEventSource{
			check:       c.Check,
			shouldError: c.ShouldError,
		})

		_, err := s.StartOrder(c.Order)
		if c.ShouldError && err == nil {
			t.Errorf("Cases[%d] FAILED: %s, expected an error.", i, c.Label)
			continue
		}

		if !c.ShouldError && err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, err)
		}
	}
}

func TestService_UpdateOrder(t *testing.T) {
	cases := []struct {
		Label       string
		Order       *model.OrderPatch
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should error if OrderID is not provided",
			Order: &model.OrderPatch{},
			Check: func(c eventsource.Command) error {
				return nil
			},
			ShouldError: true,
		},
		{
			Label: "Should correctly issue the update command",
			Order: &model.OrderPatch{
				OrderID:     "testOrderId",
				ServiceType: model.NewOptionalServiceType(model.Pickup),
				Description: optional.NewString("I'm a test!"),
			},
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.UpdateOrderCommand)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.UpdateOrderCommand{}, c)
				}
				if cmd.OrderID != "testOrderId" {
					return fmt.Errorf("Expected `%s` for OrderID, got `%s`", "testOrderId", cmd.OrderID)
				}
				if err := deep.Equal(cmd.ServiceType, model.NewOptionalServiceType(model.Pickup)); err != nil {
					return fmt.Errorf("%s", err)
				}
				if err := deep.Equal(cmd.Description, optional.NewString("I'm a test!")); err != nil {
					return fmt.Errorf("%s", err)
				}
				return nil
			},
		},
		{
			Label: "Should bubble up errors",
			Order: &model.OrderPatch{
				OrderID: "testOrderId",
			},
			Check: func(c eventsource.Command) error {
				return nil
			},
			ShouldError: true,
		},
	}

	for i, c := range cases {
		s := NewService(&mockEventSource{
			check:       c.Check,
			shouldError: c.ShouldError,
		})

		err := s.UpdateOrder(c.Order)
		if c.ShouldError && err == nil {
			t.Errorf("Cases[%d] FAILED: %s, expected an error.", i, c.Label)
			continue
		}

		if !c.ShouldError && err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, err)
		}
	}
}

func TestService_SubmitOrder(t *testing.T) {
	cases := []struct {
		Label       string
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should correctly issue the submit order command",
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.SubmitOrderCommand)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.SubmitOrderCommand{}, c)
				}
				if cmd.OrderID != "testOrderId" {
					return fmt.Errorf("Expected `%s` for OrderID, got `%s`", "testOrderId", cmd.OrderID)
				}
				if err := deep.Equal(cmd.OrderID, "testOrderId"); err != nil {
					return fmt.Errorf("%s", err)
				}
				return nil
			},
		},
		{
			Label: "Should bubble up errors",
			Check: func(c eventsource.Command) error {
				return nil
			},
			ShouldError: true,
		},
	}

	for i, c := range cases {
		s := NewService(&mockEventSource{
			check:       c.Check,
			shouldError: c.ShouldError,
		})

		err := s.SubmitOrder("testOrderId")
		if c.ShouldError && err == nil {
			t.Errorf("Cases[%d] FAILED: %s, expected an error.", i, c.Label)
			continue
		}

		if !c.ShouldError && err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, err)
		}
	}
}

func TestService_ApproveOrder(t *testing.T) {
	cases := []struct {
		Label       string
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should correctly issue the submit order command",
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.ApproveOrderCommand)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.ApproveOrderCommand{}, c)
				}
				if cmd.OrderID != "testOrderId" {
					return fmt.Errorf("Expected `%s` for OrderID, got `%s`", "testOrderId", cmd.OrderID)
				}
				if err := deep.Equal(cmd.OrderID, "testOrderId"); err != nil {
					return fmt.Errorf("%s", err)
				}
				return nil
			},
		},
		{
			Label: "Should bubble up errors",
			Check: func(c eventsource.Command) error {
				return nil
			},
			ShouldError: true,
		},
	}

	for i, c := range cases {
		s := NewService(&mockEventSource{
			check:       c.Check,
			shouldError: c.ShouldError,
		})

		err := s.ApproveOrder("testOrderId")
		if c.ShouldError && err == nil {
			t.Errorf("Cases[%d] FAILED: %s, expected an error.", i, c.Label)
			continue
		}

		if !c.ShouldError && err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, err)
		}
	}
}

func TestService_DeliverOrder(t *testing.T) {
	cases := []struct {
		Label       string
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should correctly issue the submit order command",
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.DeliverOrderCommand)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.DeliverOrderCommand{}, c)
				}
				if cmd.OrderID != "testOrderId" {
					return fmt.Errorf("Expected `%s` for OrderID, got `%s`", "testOrderId", cmd.OrderID)
				}
				if err := deep.Equal(cmd.OrderID, "testOrderId"); err != nil {
					return fmt.Errorf("%s", err)
				}
				return nil
			},
		},
		{
			Label: "Should bubble up errors",
			Check: func(c eventsource.Command) error {
				return nil
			},
			ShouldError: true,
		},
	}

	for i, c := range cases {
		s := NewService(&mockEventSource{
			check:       c.Check,
			shouldError: c.ShouldError,
		})

		err := s.DeliverOrder("testOrderId")
		if c.ShouldError && err == nil {
			t.Errorf("Cases[%d] FAILED: %s, expected an error.", i, c.Label)
			continue
		}

		if !c.ShouldError && err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, err)
		}
	}
}

type Condition func(c eventsource.Command) error

type mockEventSource struct {
	eventsource.EventSourceAPI
	check       Condition
	shouldError bool
}

func (m *mockEventSource) ProcessCommand(got eventsource.Command, a eventsource.Aggregate) error {

	if m.shouldError {
		return fmt.Errorf("Expecting an error here.")
	}

	if _, ok := a.(*Aggregate); !ok {
		return fmt.Errorf("Error: Expected %T, got %T", &Aggregate{}, a)
	}

	return m.check(got)
}
