package orderfulfillment

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	approvalEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/event"
	deliveryEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/event"
	orderEvents "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
)

type OrderFulfillmentSaga struct {
	deliverySvc delivery.ServiceAPI
	approvalSvc approval.ServiceAPI
	orderSvc    order.ServiceAPI

	OrderID         string
	Description     string
	IsDeliveryOrder bool
	Approved        bool
	Delivered       bool
}

var _ saga.SagaAPI = (*OrderFulfillmentSaga)(nil)

func New(orderSvc order.ServiceAPI, deliverySvc delivery.ServiceAPI, approvalSvc approval.ServiceAPI) *OrderFulfillmentSaga {
	return &OrderFulfillmentSaga{
		deliverySvc: deliverySvc,
		approvalSvc: approvalSvc,
		orderSvc:    orderSvc,
	}
}

func (s *OrderFulfillmentSaga) Type() string {
	return "OrderFulfillmentSaga"
}

func (s *OrderFulfillmentSaga) Version() int {
	return 1
}

func (s *OrderFulfillmentSaga) StartEvent() string {
	return "OrderStartedEvent"
}

func (s *OrderFulfillmentSaga) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OrderFulfillmentSaga) AssociationID(event eventsource.Event) (*saga.SagaAssociation, error) {
	switch d := event.Data.(type) {
	case *orderEvents.OrderStartedEvent:
		return &saga.SagaAssociation{
			ID:              d.OrderID,
			AssociationType: "OrderID",
		}, nil
	case *orderEvents.OrderDescriptionSet:
		return &saga.SagaAssociation{
			ID:              d.OrderID,
			AssociationType: "OrderID",
		}, nil
	case *orderEvents.OrderSubmitted:
		return &saga.SagaAssociation{
			ID:              d.OrderID,
			AssociationType: "OrderID",
		}, nil
	case *approvalEvents.ApprovalReceived:
		return &saga.SagaAssociation{
			ID:              strconv.Itoa(d.ApprovalID),
			AssociationType: "ApprovalID",
		}, nil
	case *deliveryEvents.DeliveryConfirmed:
		return &saga.SagaAssociation{
			ID:              strconv.Itoa(d.DeliveryID),
			AssociationType: "DeliveryID",
		}, nil
	default:
		return nil, fmt.Errorf("Unsupported event %T received: %+v", d, event)
	}
}

func (s *OrderFulfillmentSaga) HandleEvent(event eventsource.Event) (*saga.HandleEventResult, error) {
	switch d := event.Data.(type) {
	case *orderEvents.OrderStartedEvent:
		s.OrderID = d.OrderID
		s.Description = d.Description
		s.IsDeliveryOrder = d.ServiceType == model.Delivery
		return nil, nil
	case *orderEvents.OrderDescriptionSet:
		s.Description = d.Description
		return nil, nil
	case *orderEvents.OrderServiceTypeSetEvent:
		s.IsDeliveryOrder = d.ServiceType == model.Delivery
		return nil, nil
	case *orderEvents.OrderSubmitted:
		ids, err := s.startSaga(d)
		if err != nil {
			return nil, err
		}
		return &saga.HandleEventResult{AssociationIDs: ids}, nil
	case *approvalEvents.ApprovalReceived:
		return s.handleApprovalReceived(d)
	case *deliveryEvents.DeliveryConfirmed:
		return s.handleDeliveryConfirmed(d)
	default:
		return nil, fmt.Errorf("Unsupported event %T received: %+v", d, event)
	}
}

func (s *OrderFulfillmentSaga) startSaga(_ *orderEvents.OrderSubmitted) ([]*saga.SagaAssociation, error) {
	a, err := s.approvalSvc.SubmitOrderForApproval(&approval.OrderApproval{
		Description: s.Description,
	})
	if err != nil {
		return nil, err
	}
	associations := []*saga.SagaAssociation{
		{
			ID:              strconv.Itoa(a.ApprovalID),
			AssociationType: "ApprovalID",
		},
	}
	return associations, nil
}

func (s *OrderFulfillmentSaga) handleApprovalReceived(_ *approvalEvents.ApprovalReceived) (*saga.HandleEventResult, error) {

	if err := s.orderSvc.ApproveOrder(s.OrderID); err != nil {
		return nil, err
	}

	s.Approved = true

	if !s.IsDeliveryOrder {
		log.Printf("This order is not marked for delivery, skipping.")
		return nil, nil
	}

	a, err := s.deliverySvc.SubmitOrderForDelivery(&delivery.OrderDelivery{
		Description: s.Description,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Order submitted for delivery, with callback ID: %d", a.DeliveryID)

	associations := []*saga.SagaAssociation{
		{
			ID:              strconv.Itoa(a.DeliveryID),
			AssociationType: "DeliveryID",
		},
	}
	return &saga.HandleEventResult{AssociationIDs: associations}, nil
}

func (s *OrderFulfillmentSaga) handleDeliveryConfirmed(_ *deliveryEvents.DeliveryConfirmed) (*saga.HandleEventResult, error) {

	if err := s.orderSvc.DeliverOrder(s.OrderID); err != nil {
		return nil, err
	}

	s.Delivered = true

	log.Printf("Order has been delivered, fulfillment is complete!")
	return nil, nil
}
