package command

import "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

// UpdateOrderCommand allows updating multiple fields at once
type UpdateOrderCommand struct {
	OrderID     string
	Description string
	ServiceType model.ServiceType
}

func (c *UpdateOrderCommand) AggregateID() string {
	return c.OrderID
}
