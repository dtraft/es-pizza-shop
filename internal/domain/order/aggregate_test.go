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

var submitOrderCommand = &command.SubmitOrderCommand{
	OrderID: "testOrderId",
}

var approveOrderCommand = &command.ApproveOrderCommand{
	OrderID: "testOrderId",
}

var deliverOrderCommand = &command.DeliverOrderCommand{
	OrderID: "testOrderId",
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

var orderSubmittedEvent = &event.OrderSubmitted{
	OrderID: "testOrderId",
}

var orderApprovedEvent = &event.OrderApproved{
	OrderID: "testOrderId",
}

var orderDeliveredEvent = &event.OrderDelivered{
	OrderID: "testOrderId",
}

func TestAggregate_HandleCommand(t *testing.T) {

	cases := eventsourcetest.HandleCommandCases{
		{
			Label: "prevents starting orders which have already been started",
			Given: []eventsource.EventData{
				orderStartedEvent,
			},
			Command:     startOrderCommand,
			ShouldError: true,
		},
		{
			Label:   "correctly processes OrderStartedCommand",
			Command: startOrderCommand,
			Expected: []eventsource.EventData{
				orderStartedEvent,
			},
		},
		{
			Label:       "prevents UpdateOrderCommand which do not exist",
			Given:       nil,
			Command:     updateOrderCommand,
			ShouldError: true,
		},
		{
			Label: "prevents UpdateOrderCommand which are not in the started status",
			Given: []eventsource.EventData{
				orderStartedEvent,
				orderSubmittedEvent,
			},
			Command:     updateOrderCommand,
			ShouldError: true,
		},
		{
			Label: "correctly processes UpdateOrderCommand with full update",
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
			Label: "correctly processes UpdateOrderCommand with partial update",
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
			Label: "UpdateOrderCommand doesn't emit an event when no differences are detected",
			Given: []eventsource.EventData{
				orderStartedEvent,
			},
			Command:  updateOrderCommandNoUpdates,
			Expected: []eventsource.EventData{},
		},
		{
			Label:       "prevents submissions for nonexistent orders",
			Given:       nil,
			Command:     submitOrderCommand,
			ShouldError: true,
		},
		{
			Label: "prevents submissions for previously submitted orders",
			Given: []eventsource.EventData{
				orderStartedEvent,
				orderSubmittedEvent,
			},
			Command:     submitOrderCommand,
			ShouldError: true,
		},
		{
			Label: "correctly processes SubmitOrderCommand",
			Given: []eventsource.EventData{
				orderStartedEvent,
			},
			Command: submitOrderCommand,
			Expected: []eventsource.EventData{
				orderSubmittedEvent,
			},
		},
		{
			Label:       "prevents approvals for nonexistent orders",
			Given:       nil,
			Command:     approveOrderCommand,
			ShouldError: true,
		},
		{
			Label: "prevents approvals for previously approved orders",
			Given: []eventsource.EventData{
				orderStartedEvent,
				orderSubmittedEvent,
				orderApprovedEvent,
			},
			Command:     approveOrderCommand,
			ShouldError: true,
		},
		{
			Label: "correctly processes ApproveOrderCommand",
			Given: []eventsource.EventData{
				orderStartedEvent,
				orderSubmittedEvent,
			},
			Command: approveOrderCommand,
			Expected: []eventsource.EventData{
				orderApprovedEvent,
			},
		},
		{
			Label:       "prevents deliveries for nonexistent orders",
			Given:       nil,
			Command:     deliverOrderCommand,
			ShouldError: true,
		},
		{
			Label: "prevents deliveries for previously delivered orders",
			Given: []eventsource.EventData{
				orderStartedEvent,
				orderSubmittedEvent,
				orderApprovedEvent,
				orderDeliveredEvent,
			},
			Command:     deliverOrderCommand,
			ShouldError: true,
		},
		{
			Label: "correctly processes DeliverOrderCommand",
			Given: []eventsource.EventData{
				orderStartedEvent,
				orderSubmittedEvent,
				orderApprovedEvent,
			},
			Command: deliverOrderCommand,
			Expected: []eventsource.EventData{
				orderDeliveredEvent,
			},
		},
	}

	cases.Test(&Aggregate{}, t)
}

func TestAggregate_ApplyEvent(t *testing.T) {

	cases := []*eventsourcetest.ApplyEventCase{
		{
			Event: orderStartedEvent,
			Expected: &Aggregate{
				ServiceType: model.Pickup,
				Description: "Here is a description",
				Status:      model.Started,
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
				Status:      model.Started,
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
				Status:      model.Started,
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
