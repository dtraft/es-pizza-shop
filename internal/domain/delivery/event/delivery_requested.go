package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

func init() {
	eventsource.RegisterEventType(&DeliveryRequested{})
}

// DeliveryRequested fired when approval is received from vendor system
type DeliveryRequested struct {
	DeliveryID int `json:"deliveryId"`
}

func (e *DeliveryRequested) Version() int {
	return 1
}

func (e *DeliveryRequested) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}
