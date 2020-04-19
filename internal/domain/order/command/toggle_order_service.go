package command

// ToggleOrderServiceTypeCommand toggles the service type between pickup and delivery
type ToggleOrderServiceTypeCommand struct {
	OrderID string
}

func (c *ToggleOrderServiceTypeCommand) AggregateID() string {
	return c.OrderID
}
