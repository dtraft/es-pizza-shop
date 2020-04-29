package approval

import (
	"fmt"
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/command"
)

// TODO - Test Request Order

func TestService_ReceiveApproval(t *testing.T) {
	cases := []struct {
		Label       string
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should correctly issue the ReceiveApproval command",
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.ReceiveApproval)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.ReceiveApproval{}, c)
				}
				if cmd.ApprovalID != 101 {
					return fmt.Errorf("Expected `%d` for ApprovalID, got `%d`", 101, cmd.ApprovalID)
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

		err := s.ReceiveApproval(101)
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
