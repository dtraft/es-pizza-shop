package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
)

func init() {
	eventsource.RegisterEventType(&OrderStartedEvent{})
}

// OrderStartedEvent fired when an order is started
type OrderStartedEvent struct {
	OrderID     string            `json:"orderId"`
	ServiceType model.ServiceType `json:"serviceType"`
	Description string            `json:"description"`
}

func (e *OrderStartedEvent) Version() int {
	return 2
}

func (e *OrderStartedEvent) Load(data json.RawMessage, version int) error {
	switch version {
	case 1:
		v1 := OrderStartedEventV1{}
		err := json.Unmarshal(data, &v1)
		if err != nil {
			return err
		}
		e.OrderID = v1.OrderID
		e.ServiceType = model.ServiceType(v1.ServiceType)
		e.Description = v1.Description
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}

type OrderStartedEventV1 struct {
	OrderID     string `json:"orderId"`
	ServiceType int    `json:"serviceType"`
	Description string `json:"description"`
}
