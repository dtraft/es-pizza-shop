package event

import (
	"encoding/json"
	"fmt"

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
	return 1
}

func (e *OrderServiceTypeSetEvent) Load(data json.RawMessage, version int) error {
	switch version {
	case 1:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	default:
		_, eventType := eventsource.GetTypeName(e)
		return fmt.Errorf("Version %d is not supported by %s", version, eventType)
	}
	return nil
}
