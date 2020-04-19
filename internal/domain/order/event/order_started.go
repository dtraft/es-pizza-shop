package event

import (
	"encoding/json"
	"errors"
	"fmt"

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
}

func (e *OrderStartedEvent) Version() int {
	return 1
}

func (e *OrderStartedEvent) Load(data json.RawMessage, version int) error {
	switch version {
	case 1:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}

	default:
		_, eventType := eventsource.GetTypeName(e)
		return errors.New(fmt.Sprintf("Version %d is not supported by %s", version, eventType))
	}
	return nil
}
