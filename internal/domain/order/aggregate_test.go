package order

import (
	"testing"

	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/command"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
)

var startOrderCommand = &command.StartOrderCommand{
	OrderID:     "testOrderId",
	Description: "Here is a description",
	Type:        model.Pickup,
}

var toggleServiceType = &command.ToggleOrderServiceTypeCommand{
	OrderID: "testOrderId",
}

var orderStartedEvent = &event.OrderStartedEvent{
	OrderID:     "testOrderId",
	Description: "Here is a description",
	ServiceType: model.Pickup,
}

var serviceTypeSetEvent = &event.OrderServiceTypeSetEvent{
	OrderID:     "testOrderId",
	ServiceType: model.Delivery,
}

func TestAggregate_HandleCommand(t *testing.T) {

	cases := []*eventsourcetest.HandleCommandCase{
		{
			Command: startOrderCommand,
			Expected: []eventsource.EventData{
				orderStartedEvent,
			},
		},
		{
			Given: []eventsource.EventData{
				orderStartedEvent,
			},
			Command: toggleServiceType,
			Expected: []eventsource.EventData{
				serviceTypeSetEvent,
			},
		},
	}

	for i, c := range cases {
		if err := c.HandleCommand(&Aggregate{}); err != nil {
			t.Errorf("Error is cases[%d]: %s", i, err)
		}
	}
}

func TestAggregate_ApplyEvent(t *testing.T) {

	cases := []*eventsourcetest.ApplyEventCase{
		{
			Event: orderStartedEvent,
			Expected: &Aggregate{
				ServiceType: model.Pickup,
			},
		},
		{
			Given: []eventsource.EventData{
				orderStartedEvent,
			},
			Event: serviceTypeSetEvent,
			Expected: &Aggregate{
				ServiceType: model.Delivery,
			},
		},
	}

	for i, c := range cases {
		if err := c.ApplyEvent(&Aggregate{}); err != nil {
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
