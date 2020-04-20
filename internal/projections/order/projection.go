package order

import (
	"fmt"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
	. "forge.lmig.com/n1505471/pizza-shop/internal/projections/order/model"
	"forge.lmig.com/n1505471/pizza-shop/internal/projections/order/repository"
)

// Projection
type Projection struct {
	repo repository.Interface
}

func NewProjection(repo repository.Interface) *Projection {
	return &Projection{
		repo: repo,
	}
}

func (p *Projection) ApplyEvent(e es.Event) error {
	switch d := e.Data.(type) {
	case *event.OrderStartedEvent:
		return p.handleOrderStartedEvent(d)
	case *event.OrderServiceTypeSetEvent:
		return p.handleServiceTypeSetEvent(d)
	default:
		return fmt.Errorf("Unsupported event %s received in ApplyEvent handler of the Order Projection: %+v", d, e)
	}
}

/*
 * Event Handlers
 */

func (p *Projection) handleOrderStartedEvent(e *event.OrderStartedEvent) error {
	fmt.Printf("Handling projection for OrderStartedEvent: %+v", e)
	return p.repo.Save(&Order{
		OrderID:     e.OrderID,
		ServiceType: e.ServiceType,
	})
}

func (p *Projection) handleServiceTypeSetEvent(e *event.OrderServiceTypeSetEvent) error {
	fmt.Printf("Handling projection for OrderServiceTypeSetEvent: %+v", e)
	return nil
}
