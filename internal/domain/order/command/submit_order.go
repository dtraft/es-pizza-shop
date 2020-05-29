package command

// SubmitOrderCommand attempts to submit the order for fulfillment
type SubmitOrderCommand struct {
	OrderID string
}

func (c *SubmitOrderCommand) AggregateID() string {
	return c.OrderID
}
