package command

import "strconv"

// ReceiveApproval attempts to submit the order for fulfillment
type RequestApproval struct {
	ApprovalID int
}

func (c *RequestApproval) AggregateID() string {
	return strconv.Itoa(c.ApprovalID)
}
