package command

import "strconv"

// ReceiveApproval attempts to submit the order for fulfillment
type ReceiveApproval struct {
	ApprovalID int
}

func (c *ReceiveApproval) AggregateID() string {
	return strconv.Itoa(c.ApprovalID)
}
