package order

import (
	"forge.lmig.com/n1505471/pizza-shop/app/domain/order/command"
	"forge.lmig.com/n1505471/pizza-shop/app/domain/order/model"
	"forge.lmig.com/n1505471/pizza-shop/eventsource"
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

func (s *Service) StartOrder() (string, error) {
	c := &command.StartOrderCommand{
		OrderID: uuid.New().String(),
		Type:    model.Pickup,
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
