package orderfulfillment

import (
	"testing"

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
