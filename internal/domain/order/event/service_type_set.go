package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
)

func init() {
	eventsource.RegisterEventType(&OrderServiceTypeSetEvent{})
}

// OrderServiceTypeSetEvent fired when an order's service type is set
type OrderServiceTypeSetEvent struct {
	OrderID     string            `json:"orderId"`
	ServiceType model.ServiceType `json:"serviceType"`
}

func (e *OrderServiceTypeSetEvent) Version() int {
	return 2
}

func (e *OrderServiceTypeSetEvent) Load(data json.RawMessage, version int) error {
	switch version {
	case 1:
		v1 := OrderServiceTypeSetEventV1{}
		err := json.Unmarshal(data, &v1)
		if err != nil {
			return err
		}
		e.OrderID = v1.OrderID
		e.ServiceType = model.ServiceType(v1.ServiceType)
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}

type OrderServiceTypeSetEventV1 struct {
	OrderID     string `json:"orderId"`
	ServiceType int    `json:"serviceType"`
}
