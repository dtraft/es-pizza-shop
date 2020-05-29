package event

import (
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
)

func TestOrderApproved_Load(t *testing.T) {
	cases := eventsourcetest.EventLoadTestCases{
		{
			Label:   "correctly handles version 1 event",
			Version: 1,
			Event: `
				{
					"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743"
				}
			`,
			Expected: &OrderApproved{
				OrderID: "84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
			},
		},
		{
			Label:   "version 1 returns error with invalid json",
			Version: 1,
			Event: `
				{
					"orderId": 123
				}
			`,
			Expected:    &OrderApproved{},
			ShouldError: true,
		},
	}

	cases.Test(t)
}
