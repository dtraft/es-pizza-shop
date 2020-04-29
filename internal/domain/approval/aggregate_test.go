package approval

import (
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/event"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/command"
	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
)

var receiveApproval = &command.ReceiveApproval{
	ApprovalID: 101,
}

var approvalReceivedEvent = &event.ApprovalReceived{
	ApprovalID: 101,
}

func TestAggregate_HandleCommand(t *testing.T) {

	cases := eventsourcetest.HandleCommandCases{
		{
			Label: "prevents double approvals",
			Given: []eventsource.EventData{
				approvalReceivedEvent,
			},
			Command:  receiveApproval,
			Expected: nil,
		},
		{
			Label:   "correctly issues the ApprovalReceived event",
			Given:   nil,
			Command: receiveApproval,
			Expected: []eventsource.EventData{
				approvalReceivedEvent,
			},
		},
	}

	cases.Test(&Aggregate{}, t)
}

func TestAggregate_ApplyEvent(t *testing.T) {

	cases := []*eventsourcetest.ApplyEventCase{
		{
			Event: approvalReceivedEvent,
			Expected: &Aggregate{
				Approved: true,
			},
		},
	}

	for i, c := range cases {
		if err := c.TestApplyEvent(&Aggregate{}); err != nil {
			t.Errorf("Error is cases[%d]: %s", i, err)
		}
	}

}

func TestAggregate_InitAndAggregateID(t *testing.T) {
	a := &Aggregate{}
	a.Init("aggregateId")

	if diff := deep.Equal(a.AggregateID(), "aggregateId"); diff != nil {
		t.Error(diff)
	}
}
