package orderfulfillment

import (
	"fmt"
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

	approvalEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/event"
	deliveryEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/event"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"
	orderEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"

	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
)

var testSaga = &OrderFulfillmentSaga{}

func TestOrderFulfillmentSaga_Type(t *testing.T) {
	if diff := deep.Equal("OrderFulfillmentSaga", testSaga.Type()); diff != nil {
		t.Error(diff)
	}
}

func TestOrderFulfillmentSaga_Version(t *testing.T) {
	if diff := deep.Equal(1, testSaga.Version()); diff != nil {
		t.Error(diff)
	}
}

func TestOrderFulfillmentSaga_StartEvent(t *testing.T) {
	if diff := deep.Equal("OrderStartedEvent", testSaga.StartEvent()); diff != nil {
		t.Error(diff)
	}
}

func TestOrderFulfillmentSaga_Load(t *testing.T) {
	cases := eventsourcetest.SagaLoadTestCases{
		{
			Label:   "returns an error when json is invalid",
			Version: 1,
			Saga: `
				{
			`,
			Expected:    &OrderFulfillmentSaga{},
			ShouldError: true,
		},
		{
			Label:   "correctly handles version 1 event",
			Version: 1,
			Saga: `
				{
					"orderId": "testOrderId",
					"description": "test description",
					"isDeliveryOrder": true,
					"approved": true,
					"delivered": true
				}
			`,
			Expected: &OrderFulfillmentSaga{
				OrderID:         "testOrderId",
				Description:     "test description",
				IsDeliveryOrder: true,
				Approved:        true,
				Delivered:       true,
			},
		},
	}

	cases.Test(t)
}

func TestOrderFulfillmentSaga_AssociationID(t *testing.T) {
	cases := eventsourcetest.SagaAssociationIDTestCases{
		{
			Label: "handles OrderStartedEvent",
			Saga:  &OrderFulfillmentSaga{},
			Event: eventsource.Event{Data: &orderEvents.OrderStartedEvent{
				OrderID: "orderID",
			}},
			Expected: &saga.SagaAssociation{
				ID:              "orderID",
				AssociationType: "OrderID",
			},
		},
		{
			Label: "handles OrderDescriptionSet",
			Saga:  &OrderFulfillmentSaga{},
			Event: eventsource.Event{Data: &orderEvents.OrderDescriptionSet{
				OrderID: "orderID",
			}},
			Expected: &saga.SagaAssociation{
				ID:              "orderID",
				AssociationType: "OrderID",
			},
		},
		{
			Label: "handles OrderSubmitted",
			Saga:  &OrderFulfillmentSaga{},
			Event: eventsource.Event{Data: &orderEvents.OrderSubmitted{
				OrderID: "orderID",
			}},
			Expected: &saga.SagaAssociation{
				ID:              "orderID",
				AssociationType: "OrderID",
			},
		},
		{
			Label: "handles ApprovalReceived",
			Saga:  &OrderFulfillmentSaga{},
			Event: eventsource.Event{Data: &approvalEvents.ApprovalReceived{
				ApprovalID: 1,
			}},
			Expected: &saga.SagaAssociation{
				ID:              "1",
				AssociationType: "ApprovalID",
			},
		},
		{
			Label: "handles ApprovalReceived",
			Saga:  &OrderFulfillmentSaga{},
			Event: eventsource.Event{Data: &deliveryEvents.DeliveryConfirmed{
				DeliveryID: 1,
			}},
			Expected: &saga.SagaAssociation{
				ID:              "1",
				AssociationType: "DeliveryID",
			},
		},
		{
			Label:       "returns error for unsupported event types",
			Saga:        &OrderFulfillmentSaga{},
			Event:       eventsource.Event{Data: nil},
			ShouldError: true,
		},
	}

	cases.Test(t)
}

var orderStartedEvent = eventsource.Event{Data: &orderEvents.OrderStartedEvent{
	OrderID:     "orderID",
	Description: "test description",
	ServiceType: model.Pickup,
}}

var descriptionSetEvent = eventsource.Event{Data: &orderEvents.OrderDescriptionSet{
	OrderID:     "orderID",
	Description: "test description, again",
}}

var serviceTypeSetEvent = eventsource.Event{Data: &orderEvents.OrderServiceTypeSetEvent{
	OrderID:     "orderID",
	ServiceType: model.Delivery,
}}

var submittedEvent = eventsource.Event{Data: &orderEvents.OrderSubmitted{
	OrderID: "orderID",
}}

var approvalReceived = eventsource.Event{Data: &approvalEvents.ApprovalReceived{
	ApprovalID: 1,
}}

var deliveryConfirmed = eventsource.Event{Data: &deliveryEvents.DeliveryConfirmed{
	DeliveryID: 2,
}}

func TestOrderFulfillmentSaga_HandleEvent(t *testing.T) {
	cases := eventsourcetest.SagaHandleEventTestCases{
		{
			Label: "handles OrderStartedEvent",
			Event: orderStartedEvent,
			ExpectedSaga: &OrderFulfillmentSaga{
				OrderID:     "orderID",
				Description: "test description",
			},
		},
		{
			Label: "handles OrderDescriptionSet",
			Given: []eventsource.Event{
				orderStartedEvent,
			},
			Event: descriptionSetEvent,
			ExpectedSaga: &OrderFulfillmentSaga{
				OrderID:     "orderID",
				Description: "test description, again",
			},
		},
		{
			Label: "handles OrderServiceTypeSetEvent",
			Given: []eventsource.Event{
				orderStartedEvent,
			},
			Event: serviceTypeSetEvent,
			ExpectedSaga: &OrderFulfillmentSaga{
				OrderID:         "orderID",
				Description:     "test description",
				IsDeliveryOrder: true,
			},
		},
		{
			Label: "handles OrderSubmitted correctly",
			Saga: New(&mockOrderSvc{}, &mockDeliverySvc{}, &mockApprovalSvc{Expected: &approval.OrderApproval{
				Description: "test description",
			}}),
			Given: []eventsource.Event{
				orderStartedEvent,
				serviceTypeSetEvent,
			},
			Event: submittedEvent,
			ExpectedSaga: &OrderFulfillmentSaga{
				OrderID:         "orderID",
				Description:     "test description",
				IsDeliveryOrder: true,
			},
			ExpectedResult: &saga.HandleEventResult{AssociationIDs: []*saga.SagaAssociation{{
				AssociationType: "ApprovalID",
				ID:              "1",
			}}},
		},
		{
			Label: "forwards errors from approval service on OrderSubmitted",
			Saga:  New(&mockOrderSvc{}, &mockDeliverySvc{}, &mockApprovalSvc{ShouldError: true}),
			Given: []eventsource.Event{
				orderStartedEvent,
				serviceTypeSetEvent,
			},
			Event:         submittedEvent,
			ShouldError:   true,
			ExpectedError: fmt.Errorf("Error in SubmitOrderForApproval"),
		},
		{
			Label: "handles ApprovalReceived correctly",
			Saga: New(
				&mockOrderSvc{Expected: "orderID"},
				&mockDeliverySvc{Expected: &delivery.OrderDelivery{
					Description: "test description",
				}},
				&mockApprovalSvc{},
			),
			Given: []eventsource.Event{
				orderStartedEvent,
				serviceTypeSetEvent,
				submittedEvent,
			},
			Event: approvalReceived,
			ExpectedSaga: &OrderFulfillmentSaga{
				OrderID:         "orderID",
				Description:     "test description",
				IsDeliveryOrder: true,
				Approved:        true,
			},
			ExpectedResult: &saga.HandleEventResult{AssociationIDs: []*saga.SagaAssociation{{
				AssociationType: "DeliveryID",
				ID:              "2",
			}}},
		},
		{
			Label: "forwards errors from order service on ApprovalReceived",
			Saga:  New(&mockOrderSvc{ShouldError: true}, &mockDeliverySvc{}, &mockApprovalSvc{}),
			Given: []eventsource.Event{
				orderStartedEvent,
				serviceTypeSetEvent,
				submittedEvent,
			},
			Event:         approvalReceived,
			ShouldError:   true,
			ExpectedError: fmt.Errorf("Error in ApproveOrder"),
		},
		{
			Label: "skips orders not intended for delivery on ApprovalReceived",
			Saga:  New(&mockOrderSvc{}, &mockDeliverySvc{}, &mockApprovalSvc{}),
			Given: []eventsource.Event{
				orderStartedEvent,
				submittedEvent,
			},
			Event:          approvalReceived,
			ExpectedResult: nil,
			ExpectedSaga: &OrderFulfillmentSaga{
				OrderID:         "orderID",
				Description:     "test description",
				IsDeliveryOrder: false,
				Approved:        true,
			},
		},
		{
			Label: "forwards errors from delivery service on ApprovalReceived",
			Saga:  New(&mockOrderSvc{}, &mockDeliverySvc{ShouldError: true}, &mockApprovalSvc{}),
			Given: []eventsource.Event{
				orderStartedEvent,
				serviceTypeSetEvent,
				submittedEvent,
			},
			Event:         approvalReceived,
			ShouldError:   true,
			ExpectedError: fmt.Errorf("Error in SubmitOrderForDelivery"),
		},
		{
			Label: "handles DeliveryConfirmed correctly",
			Saga: New(
				&mockOrderSvc{Expected: "orderID"},
				&mockDeliverySvc{},
				&mockApprovalSvc{},
			),
			Given: []eventsource.Event{
				orderStartedEvent,
				serviceTypeSetEvent,
				submittedEvent,
				approvalReceived,
			},
			Event: deliveryConfirmed,
			ExpectedSaga: &OrderFulfillmentSaga{
				OrderID:         "orderID",
				Description:     "test description",
				IsDeliveryOrder: true,
				Approved:        true,
				Delivered:       true,
			},
			ExpectedResult: nil,
		},
		{
			Label: "forwards errors from order service on DeliveryConfirmed",
			Saga:  New(&mockOrderSvc{ShouldError: true}, &mockDeliverySvc{}, &mockApprovalSvc{}),
			Given: []eventsource.Event{
				orderStartedEvent,
				serviceTypeSetEvent,
				submittedEvent,
				approvalReceived,
			},
			Event:         deliveryConfirmed,
			ShouldError:   true,
			ExpectedError: fmt.Errorf("Error in DeliverOrder"),
		},
		{
			Label:       "returns error for unsupported event types",
			Saga:        &OrderFulfillmentSaga{},
			Event:       eventsource.Event{Data: nil},
			ShouldError: true,
		},
	}

	cases.Test(t)
}

type mockOrderSvc struct {
	order.ServiceAPI
	Expected    interface{}
	ShouldError bool
}

func (m *mockOrderSvc) ApproveOrder(orderId string) error {
	if m.ShouldError {
		return fmt.Errorf("Error in ApproveOrder")
	}

	if diff := deep.Equal(orderId, m.Expected); m.Expected != nil && diff != nil {
		return fmt.Errorf("OrderID in OrderService does not match expected, details: %s", diff)
	}

	return nil
}

func (m *mockOrderSvc) DeliverOrder(orderId string) error {
	if m.ShouldError {
		return fmt.Errorf("Error in DeliverOrder")
	}

	if diff := deep.Equal(orderId, m.Expected); m.Expected != nil && diff != nil {
		return fmt.Errorf("OrderID in OrderService does not match expected, details: %s", diff)
	}

	return nil
}

type mockApprovalSvc struct {
	approval.ServiceAPI
	Expected    *approval.OrderApproval
	ShouldError bool
}

func (m *mockApprovalSvc) SubmitOrderForApproval(a *approval.OrderApproval) (*approval.OrderApproval, error) {
	if m.ShouldError {
		return nil, fmt.Errorf("Error in SubmitOrderForApproval")
	}

	if diff := deep.Equal(a, m.Expected); m.Expected != nil && diff != nil {
		return nil, fmt.Errorf("OrderApproval does not match expected, details: %s", diff)
	}

	a.ApprovalID = 1

	return a, nil
}

type mockDeliverySvc struct {
	delivery.ServiceAPI
	Expected    interface{}
	ShouldError bool
}

func (m *mockDeliverySvc) SubmitOrderForDelivery(d *delivery.OrderDelivery) (*delivery.OrderDelivery, error) {
	if m.ShouldError {
		return nil, fmt.Errorf("Error in SubmitOrderForDelivery")
	}

	if diff := deep.Equal(d, m.Expected); m.Expected != nil && diff != nil {
		return nil, fmt.Errorf("OrderApproval does not match expected, details: %s", diff)
	}

	d.DeliveryID = 2

	return d, nil
}
