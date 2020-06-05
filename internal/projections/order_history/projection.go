package order

import (
	"encoding/json"
	"fmt"
	"log"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
	. "forge.lmig.com/n1505471/pizza-shop/internal/projections/order_history/model"
	"forge.lmig.com/n1505471/pizza-shop/internal/projections/order_history/repository"
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
	log.Printf("Handling %s.\n", e.EventType)
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
		log.Printf("Unsupported event %T received in handler of the Order Projection: %+v", d, e)
		return nil
	}
}

/*
 * Event Handlers
 */

func (p *Projection) handleOrderStartedEvent(d *event.OrderStartedEvent, e es.Event) error {
	serviceType, _ := json.Marshal(d.ServiceType)
	return p.repo.Save(&OrderHistoryRecord{
		OrderID:         d.OrderID,
		Description:     fmt.Sprintf("Order started with description: %s and serviceType: %s", d.Description, string(serviceType)),
		TransactionDate: &e.Timestamp,
	})
}

func (p *Projection) handleServiceTypeSetEvent(d *event.OrderServiceTypeSetEvent, e es.Event) error {
	serviceType, _ := json.Marshal(d.ServiceType)
	return p.repo.Save(&OrderHistoryRecord{
		OrderID:         d.OrderID,
		Description:     fmt.Sprintf("ServiceType updated to: %s", string(serviceType)),
		TransactionDate: &e.Timestamp,
	})
}

func (p *Projection) handleDescriptionSetEvent(d *event.OrderDescriptionSet, e es.Event) error {
	return p.repo.Save(&OrderHistoryRecord{
		OrderID:         d.OrderID,
		Description:     fmt.Sprintf("Description updated to: %s", d.Description),
		TransactionDate: &e.Timestamp,
	})
}

func (p *Projection) handleSubmittedEvent(d *event.OrderSubmitted, e es.Event) error {
	return p.repo.Save(&OrderHistoryRecord{
		OrderID:         d.OrderID,
		Description:     "Order Submitted.",
		TransactionDate: &e.Timestamp,
	})
}

func (p *Projection) handleApprovedEvent(d *event.OrderApproved, e es.Event) error {

	return p.repo.Save(&OrderHistoryRecord{
		OrderID:         d.OrderID,
		Description:     "Order Approved.",
		TransactionDate: &e.Timestamp,
	})
}

func (p *Projection) handleDeliveredEvent(d *event.OrderDelivered, e es.Event) error {
	return p.repo.Save(&OrderHistoryRecord{
		OrderID:         d.OrderID,
		Description:     "Order Delivered.",
		TransactionDate: &e.Timestamp,
	})
}
