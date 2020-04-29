package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

func init() {
	eventsource.RegisterEventType(&OrderSubmitted{})
}

// OrderSubmitted fired when the order is submitted
type OrderSubmitted struct {
	OrderID string `json:"orderId"`
}

func (e *OrderSubmitted) Version() int {
	return 1
}

func (e *OrderSubmitted) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}
