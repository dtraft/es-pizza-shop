package order

import (
	"fmt"
	"log"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

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
	case *event.OrderSubmitted:
		return p.handleSubmittedEvent(d, e)
	case *event.OrderApproved:
		return p.handleApprovedEvent(d, e)
	case *event.OrderDelivered:
		return p.handleDeliveredEvent(d, e)
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
		Status:      model.Started,
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

func (p *Projection) handleSubmittedEvent(d *event.OrderSubmitted, e es.Event) error {

	return p.repo.Patch(d.OrderID, &Order{
		Status:    model.Submitted,
		UpdatedAt: &e.Timestamp,
	})
}

func (p *Projection) handleApprovedEvent(d *event.OrderApproved, e es.Event) error {

	return p.repo.Patch(d.OrderID, &Order{
		Status:    model.Approved,
		UpdatedAt: &e.Timestamp,
	})
}

func (p *Projection) handleDeliveredEvent(d *event.OrderDelivered, e es.Event) error {

	return p.repo.Patch(d.OrderID, &Order{
		Status:    model.Delivered,
		UpdatedAt: &e.Timestamp,
	})
}
