package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

func init() {
	eventsource.RegisterEventType(&OrderDelivered{})
}

// OrderDelivered fired when the order is submitted
type OrderDelivered struct {
	OrderID string `json:"orderId"`
}

func (e *OrderDelivered) Version() int {
	return 1
}

func (e *OrderDelivered) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}
