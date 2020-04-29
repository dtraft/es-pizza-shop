package event

import (
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/eventsourcetest"
)

func TestApprovalRequested_Load(t *testing.T) {
	cases := eventsourcetest.LoadCases{
		{
			Label:   "correctly handles version 1 event",
			Version: 1,
			Event: `
				{
					"approvalId": 101
				}
			`,
			Expected: &ApprovalRequested{
				ApprovalID: 101,
			},
		},
		{
			Label:   "version 1 returns error with invalid json",
			Version: 1,
			Event: `
				{
					"approvalID":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743"
				}
			`,
			Expected:    &ApprovalRequested{},
			ShouldError: true,
		},
	}

	cases.Test(t)
}
