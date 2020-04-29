package event

import (
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
)

func TestApprovalReceived_Load(t *testing.T) {
	cases := eventsourcetest.LoadCases{
		{
			Label:   "correctly handles version 1 event",
			Version: 1,
			Event: `
				{
					"deliveryId": 101
				}
			`,
			Expected: &DeliveryConfirmed{
				DeliveryID: 101,
			},
		},
		{
			Label:   "version 1 returns error with invalid json",
			Version: 1,
			Event: `
				{
					"deliveryId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743"
				}
			`,
			Expected:    &DeliveryConfirmed{},
			ShouldError: true,
		},
	}

	cases.Test(t)
}