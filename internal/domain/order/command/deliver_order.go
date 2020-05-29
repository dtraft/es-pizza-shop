package command

// ApproveOrderCommand attempts to submit the order for fulfillment
type DeliverOrderCommand struct {
	OrderID string
}

func (c *DeliverOrderCommand) AggregateID() string {
	return c.OrderID
}
