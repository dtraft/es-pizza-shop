package event

import (
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
)

func TestOrderServiceTypeSetEvent_Load(t *testing.T) {
	cases := eventsourcetest.EventLoadTestCases{
		{
			Label:   "correctly handles version 1 event",
			Version: 1,
			Event: `
				{
					"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					"serviceType": 1
				}
			`,
			Expected: &OrderServiceTypeSetEvent{
				OrderID:     "84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
				ServiceType: model.Pickup,
			},
		},
		{
			Label:   "version 1 returns error with invalid json",
			Version: 1,
			Event: `
				{
					"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					"serviceType": "test"
				}
			`,
			Expected:    &OrderServiceTypeSetEvent{},
			ShouldError: true,
		},
		{
			Label:   "correctly handles version 2 event",
			Version: 2,
			Event: `
				{
					"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					"serviceType": "Pickup",
					"description": "a test!"
				}
			`,
			Expected: &OrderServiceTypeSetEvent{
				OrderID:     "84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
				ServiceType: model.Pickup,
			},
		},
		{
			Label:   "version 2 returns error with invalid json",
			Version: 2,
			Event: `
				{
					"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					"serviceType": "test"
				}
			`,
			Expected:    &OrderServiceTypeSetEvent{},
			ShouldError: true,
		},
	}

	cases.Test(t)
}
