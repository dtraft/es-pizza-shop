package order

import (
	"fmt"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/command"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	"github.com/google/uuid"
)

type ServiceAPI interface {
	StartOrder(order *model.Order) (string, error)
	UpdateOrder(order *model.OrderPatch) error
}

type Service struct {
	eventSource eventsource.EventSourceAPI
}

func (s *Service) processCommand(c eventsource.Command) error {
	return s.eventSource.ProcessCommand(c, &Aggregate{})
}

func NewService(eventSource eventsource.EventSourceAPI) *Service {
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
		ServiceType: order.ServiceType,
		Description: order.Description,
	}

	if err := s.processCommand(c); err != nil {
		return "", err
	}

	return c.OrderID, nil
}

func (s *Service) UpdateOrder(order *model.OrderPatch) error {
	if order.OrderID == "" {
		return fmt.Errorf("OrderID must be provided to update operation.")
	}

	c := &command.UpdateOrderCommand{
		OrderID:     order.OrderID,
		ServiceType: order.ServiceType,
		Description: order.Description,
	}

	if err := s.processCommand(c); err != nil {
		return err
	}

	return nil
}

var _ ServiceAPI = (*Service)(nil)
