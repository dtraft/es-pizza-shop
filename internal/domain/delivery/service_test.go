package delivery

import (
	"fmt"
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/command"
)

func TestService_ReceiveDeliveryNotification(t *testing.T) {
	cases := []struct {
		Label       string
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should correctly issue the DeliveryConfirmed command",
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.ConfirmDelivery)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.ConfirmDelivery{}, c)
				}
				if cmd.DeliveryID != 101 {
					return fmt.Errorf("Expected `%d` for DeliveryID, got `%d`", 101, cmd.DeliveryID)
				}
				return nil
			},
		},
	}

	for i, c := range cases {
		s := NewService(&mockEventSource{
			check:       c.Check,
			shouldError: c.ShouldError,
		})

		err := s.ReceiveDeliveryNotification(101)
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
