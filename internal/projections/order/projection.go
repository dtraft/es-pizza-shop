package order

import (
	"fmt"
	"log"

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

func (p *Projection) HandleEvent(e es.Event) error {
	log.Printf("Handling %T.\n", e)
	switch d := e.Data.(type) {
	case *event.OrderStartedEvent:
		return p.handleOrderStartedEvent(d, e)
	case *event.OrderServiceTypeSetEvent:
		return p.handleServiceTypeSetEvent(d, e)
	case *event.OrderDescriptionSet:
		return p.handleDescriptionSetEvent(d, e)
	default:
		return fmt.Errorf("Unsupported event %T received in TestApplyEvent handler of the Order Projection: %+v", d, e)
	}
}

/*
 * Event Handlers
 */

func (p *Projection) handleOrderStartedEvent(d *event.OrderStartedEvent, e es.Event) error {
	return p.repo.Save(&Order{
		OrderID:     d.OrderID,
		ServiceType: d.ServiceType,
		Description: d.Description,
		CreatedAt:   &e.Timestamp,
		UpdatedAt:   &e.Timestamp,
	})
}

func (p *Projection) handleServiceTypeSetEvent(d *event.OrderServiceTypeSetEvent, e es.Event) error {
	return p.repo.Patch(d.OrderID, &Order{
		ServiceType: d.ServiceType,
		UpdatedAt:   &e.Timestamp,
	})
}

func (p *Projection) handleDescriptionSetEvent(d *event.OrderDescriptionSet, e es.Event) error {

	return p.repo.Patch(d.OrderID, &Order{
		Description: d.Description,
		UpdatedAt:   &e.Timestamp,
	})
}
