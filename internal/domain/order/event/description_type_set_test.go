package event

import (
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
)

func TestOrderDescriptionSet_Load(t *testing.T) {
	cases := eventsourcetest.EventLoadTestCases{
		{
			Label:   "correctly handles version 1 event",
			Version: 1,
			Event: `
				{
					"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					"description": "I'm a test!"
				}
			`,
			Expected: &OrderDescriptionSet{
				OrderID:     "84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
				Description: "I'm a test!",
			},
		},
		{
			Label:   "version 1 returns error with invalid json",
			Version: 1,
			Event: `
				{
					"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					"description": 2
				}
			`,
			Expected:    &OrderDescriptionSet{},
			ShouldError: true,
		},
	}

	cases.Test(t)
}
