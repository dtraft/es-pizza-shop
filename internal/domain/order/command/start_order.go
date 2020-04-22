package command

import "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

// StartOrderCommand starts an order
type StartOrderCommand struct {
	OrderID     string
	Type        model.ServiceType
	Description string
}

func (c *StartOrderCommand) AggregateID() string {
	return c.OrderID
}
