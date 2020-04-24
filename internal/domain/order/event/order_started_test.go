package event

import (
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
)

func TestOrderStartedEvent_Load(t *testing.T) {
	cases := eventsourcetest.LoadCases{
		{
			Label:   "correctly handles version 1 event",
			Version: 1,
			Event: `
				{
					"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					"serviceType": 1,
					"description": "a test!"
				}
			`,
			Expected: &OrderStartedEvent{
				OrderID:     "84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
				ServiceType: model.Pickup,
				Description: "a test!",
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
			Expected:    &OrderStartedEvent{},
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
			Expected: &OrderStartedEvent{
				OrderID:     "84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
				ServiceType: model.Pickup,
				Description: "a test!",
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
			Expected:    &OrderStartedEvent{},
			ShouldError: true,
		},
	}

	cases.Test(t)
}
