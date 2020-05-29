package order

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	. "forge.lmig.com/n1505471/pizza-shop/internal/projections/order/model"
	"github.com/go-test/deep"
)

var p = &Projection{}

var orderAgg = &order.Aggregate{}

// SETUP
func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

var startedEvent = eventsource.NewEvent(orderAgg, &event.OrderStartedEvent{
	OrderID:     "testOrderId",
	Description: "test desc",
	ServiceType: model.Pickup,
})

var serviceTypeSetEvent = eventsource.NewEvent(orderAgg, &event.OrderServiceTypeSetEvent{
	OrderID:     "testOrderId",
	ServiceType: model.Delivery,
})

var descriptionSetEvent = eventsource.NewEvent(orderAgg, &event.OrderDescriptionSet{
	OrderID:     "testOrderId",
	Description: "I'm a test!",
})

var submittedEvent = eventsource.NewEvent(orderAgg, &event.OrderSubmitted{
	OrderID: "testOrderId",
})

var approvedEvent = eventsource.NewEvent(orderAgg, &event.OrderApproved{
	OrderID: "testOrderId",
})

var deliveredEvent = eventsource.NewEvent(orderAgg, &event.OrderDelivered{
	OrderID: "testOrderId",
})

func TestProjection_ApplyEvent(t *testing.T) {
	cases := []struct {
		Event    eventsource.Event
		Expected *Order
	}{
		{
			Event: startedEvent,
			Expected: &Order{
				OrderID:     "testOrderId",
				Description: "test desc",
				ServiceType: model.Pickup,
				Status:      model.Started,
				CreatedAt:   &startedEvent.Timestamp,
				UpdatedAt:   &startedEvent.Timestamp,
			},
		},
		{

			Event: serviceTypeSetEvent,
			Expected: &Order{
				OrderID:     "testOrderId",
				ServiceType: model.Delivery,
				UpdatedAt:   &serviceTypeSetEvent.Timestamp,
			},
		},
		{
			Event: descriptionSetEvent,
			Expected: &Order{
				OrderID:     "testOrderId",
				Description: "I'm a test!",
				UpdatedAt:   &descriptionSetEvent.Timestamp,
			},
		},
		{
			Event: submittedEvent,
			Expected: &Order{
				OrderID:   "testOrderId",
				Status:    model.Submitted,
				UpdatedAt: &submittedEvent.Timestamp,
			},
		},
		{
			Event: approvedEvent,
			Expected: &Order{
				OrderID:   "testOrderId",
				Status:    model.Approved,
				UpdatedAt: &approvedEvent.Timestamp,
			},
		},
		{
			Event: deliveredEvent,
			Expected: &Order{
				OrderID:   "testOrderId",
				Status:    model.Delivered,
				UpdatedAt: &deliveredEvent.Timestamp,
			},
		},
	}

	for i, c := range cases {
		p.repo = &mockRepo{
			expected: c.Expected,
		}
		if err := p.HandleEvent(c.Event); err != nil {
			t.Errorf("Cases[%d]: %s", i, err)
		}
	}
}

type mockRepo struct {
	expected *Order
}

func (m *mockRepo) Save(got *Order) error {
	if diff := deep.Equal(got, m.expected); diff != nil {
		return fmt.Errorf("%s", diff)
	}
	return nil
}

func (m *mockRepo) Patch(orderID string, got *Order) error {
	if orderID != m.expected.OrderID {
		return fmt.Errorf("Expected %s, got %s for OrderID in Patch operation.", m.expected.OrderID, orderID)
	}

	// Reset OrderId for comparison testing
	m.expected.OrderID = ""

	if diff := deep.Equal(got, m.expected); diff != nil {
		return fmt.Errorf("%s", diff)
	}
	return nil
}
