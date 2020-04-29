package command

// ReceiveApproval attempts to submit the order for fulfillment
type ConfirmDelivery struct {
	DeliveryID int
}

func (c *ConfirmDelivery) AggregateID() string {
	return string(c.DeliveryID)
}
