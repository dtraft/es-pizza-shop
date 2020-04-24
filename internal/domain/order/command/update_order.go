package command

import (
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	"github.com/markphelps/optional"
)

// UpdateOrderCommand allows updating multiple fields at once
type UpdateOrderCommand struct {
	OrderID     string
	Description optional.String
	ServiceType model.OptionalServiceType
}

func (c *UpdateOrderCommand) AggregateID() string {
	return c.OrderID
}
