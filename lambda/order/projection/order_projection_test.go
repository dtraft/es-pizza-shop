package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	"github.com/go-test/deep"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource"
	"github.com/aws/aws-lambda-go/events"
)

// SETUP
func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestHandleEvent(t *testing.T) {

	cases := []struct {
		Record   events.SNSEventRecord
		Expected []es.Event
	}{
		{
			Record: events.SNSEventRecord{
				SNS: events.SNSEntity{
					MessageAttributes: map[string]interface{}{
						"eventType": map[string]interface{}{
							"Type":  "String",
							"Value": "OrderStartedEvent",
						},
					},
					Message: `
						{
							"eventId":"6c4539e3-ae1b-44f0-bfc2-4d7531893136",
							"aggregateId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
							"aggregateType":"OrderAggregate",
							"eventVersion":1,
							"eventType":"OrderStartedEvent",
							"eventTimestamp":"2020-04-19T19:45:11.475995951Z",
							"eventData":{
								"orderId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
								"serviceType":1
							}
						}
					`,
				},
			},
			Expected: []es.Event{
				{
					EventID:          "6c4539e3-ae1b-44f0-bfc2-4d7531893136",
					AggregateID:      "84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					AggregateType:    "OrderAggregate",
					EventTypeVersion: 1,
					EventType:        "OrderStartedEvent",
					Timestamp:        time.Date(2020, 04, 19, 19, 45, 11, 475995951, time.UTC),
					Data: &event.OrderStartedEvent{
						OrderID:     "84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
						ServiceType: model.Pickup,
					},
				},
			},
		},
	}

	for _, c := range cases {
		projection = &mockProjection{
			Expected: c.Expected,
		}
		if err := handleEvent(c.Record); err != nil {
			t.Error(err)
		}
	}

}

type mockProjection struct {
	Expected []es.Event
	index    int
}

func (m *mockProjection) HandleEvent(event es.Event) error {
	expected := m.Expected[m.index]
	m.index++

	if diff := deep.Equal(expected, event); diff != nil {
		return fmt.Errorf("%s", diff)
	}
	return nil
}
