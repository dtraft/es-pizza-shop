package order

import (
	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/command"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	"github.com/google/uuid"
)

type Service struct {
	eventSource eventsource.EventSource
}

func (s *Service) processCommand(c eventsource.Command) error {
	return s.eventSource.ProcessCommand(c, &Aggregate{})
}

func NewService(eventSource eventsource.EventSource) *Service {
	return &Service{
		eventSource: eventSource,
	}
}

func (s *Service) StartOrder(order *model.Order) (string, error) {
	if order.OrderID == "" {
		order.OrderID = uuid.New().String()
	}
	c := &command.StartOrderCommand{
		OrderID:     order.OrderID,
		Type:        order.ServiceType,
		Description: order.Description,
	}

	if err := s.processCommand(c); err != nil {
		return "", err
	}

	return c.OrderID, nil
}

func (s *Service) ToggleOrderServiceType(orderID string) error {
	c := &command.ToggleOrderServiceTypeCommand{
		OrderID: orderID,
	}

	if err := s.processCommand(c); err != nil {
		return err
	}

	return nil
}
