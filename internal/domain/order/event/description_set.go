package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

func init() {
	eventsource.RegisterEventType(&OrderDescriptionSet{})
}

// OrderServiceTypeSetEvent fired when an order's service type is set
type OrderDescriptionSet struct {
	OrderID     string `json:"orderId"`
	Description string `json:"description"`
}

func (e *OrderDescriptionSet) Version() int {
	return 1
}

func (e *OrderDescriptionSet) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}
