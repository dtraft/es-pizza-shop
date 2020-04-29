package delivery

import (
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/command"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/event"
	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
)

var confirmDelivery = &command.ConfirmDelivery{
	DeliveryID: 101,
}

var deliveryConfirmedEvent = &event.DeliveryConfirmed{
	DeliveryID: 101,
}

func TestAggregate_HandleCommand(t *testing.T) {

	cases := eventsourcetest.HandleCommandCases{
		{
			Label: "prevents double deliveries",
			Given: []eventsource.EventData{
				deliveryConfirmedEvent,
			},
			Command:  confirmDelivery,
			Expected: nil,
		},
		{
			Label:   "correctly issues the DeliveryConfirmed event",
			Given:   nil,
			Command: confirmDelivery,
			Expected: []eventsource.EventData{
				deliveryConfirmedEvent,
			},
		},
	}

	cases.Test(&Aggregate{}, t)
}

func TestAggregate_ApplyEvent(t *testing.T) {

	cases := []*eventsourcetest.ApplyEventCase{
		{
			Event: deliveryConfirmedEvent,
			Expected: &Aggregate{
				Delivered: true,
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
