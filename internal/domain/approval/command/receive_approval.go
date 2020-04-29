package command

// ReceiveApproval attempts to submit the order for fulfillment
type ReceiveApproval struct {
	ApprovalID int
}

func (c *ReceiveApproval) AggregateID() string {
	return string(c.ApprovalID)
}
