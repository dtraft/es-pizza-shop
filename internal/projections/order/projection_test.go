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

func TestProjection_ApplyEvent(t *testing.T) {
	cases := []struct {
		Event    eventsource.EventData
		Expected *Order
	}{
		{
			Event: &event.OrderStartedEvent{
				OrderID:     "testOrderId",
				Description: "test desc",
				ServiceType: model.Pickup,
			},
			Expected: &Order{
				OrderID:     "testOrderId",
				Description: "test desc",
				ServiceType: model.Pickup,
			},
		},
		{
			Event: &event.OrderServiceTypeSetEvent{
				OrderID:     "testOrderId",
				ServiceType: model.Delivery,
			},
			Expected: &Order{
				OrderID:     "testOrderId",
				ServiceType: model.Delivery,
			},
		},
	}

	for i, c := range cases {
		p.repo = &mockRepo{
			expected: c.Expected,
		}
		event := eventsource.NewEvent(orderAgg, c.Event)
		if err := p.HandleEvent(event); err != nil {
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
