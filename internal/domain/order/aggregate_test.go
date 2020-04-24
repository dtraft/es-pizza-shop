package order

import (
	"testing"

	"github.com/markphelps/optional"

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
	ServiceType: model.Pickup,
}

var updateOrderCommand = &command.UpdateOrderCommand{
	OrderID:     "testOrderId",
	Description: optional.NewString("Here is a NEW description"),
	ServiceType: model.NewOptionalServiceType(model.Delivery),
}

var updateOrderCommandNoUpdates = &command.UpdateOrderCommand{
	OrderID:     "testOrderId",
	Description: optional.NewString("Here is a description"),
	ServiceType: model.OptionalServiceType{},
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

var descriptionSetEvent = &event.OrderDescriptionSet{
	OrderID:     "testOrderId",
	Description: "Here is a NEW description",
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
			Command: updateOrderCommand,
			Expected: []eventsource.EventData{
				serviceTypeSetEvent,
				descriptionSetEvent,
			},
		},
		{
			Given: []eventsource.EventData{
				orderStartedEvent,
				serviceTypeSetEvent,
			},
			Command: updateOrderCommand,
			Expected: []eventsource.EventData{
				descriptionSetEvent,
			},
		},
		{
			Given: []eventsource.EventData{
				orderStartedEvent,
			},
			Command:  updateOrderCommandNoUpdates,
			Expected: []eventsource.EventData{},
		},
	}

	for i, c := range cases {
		if err := c.TestHandleCommand(&Aggregate{}); err != nil {
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
				Description: "Here is a description",
			},
		},
		{
			Given: []eventsource.EventData{
				orderStartedEvent,
			},
			Event: serviceTypeSetEvent,
			Expected: &Aggregate{
				ServiceType: model.Delivery,
				Description: "Here is a description",
			},
		},
		{
			Given: []eventsource.EventData{
				orderStartedEvent,
			},
			Event: descriptionSetEvent,
			Expected: &Aggregate{
				ServiceType: model.Pickup,
				Description: "Here is a NEW description",
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
