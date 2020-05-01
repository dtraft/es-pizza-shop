package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

func init() {
	eventsource.RegisterEventType(&DeliveryConfirmed{})
}

// DeliveryConfirmed fired when approval is received from vendor system
type DeliveryConfirmed struct {
	DeliveryID int `json:"deliveryId"`
}

func (e *DeliveryConfirmed) Version() int {
	return 1
}

func (e *DeliveryConfirmed) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}
