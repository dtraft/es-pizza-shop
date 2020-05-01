package command

// ApproveOrderCommand attempts to submit the order for fulfillment
type ApproveOrderCommand struct {
	OrderID string
}

func (c *ApproveOrderCommand) AggregateID() string {
	return c.OrderID
}
