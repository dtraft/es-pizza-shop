package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleEvent(t *testing.T) {

	cases := []struct {
		Record events.SNSEventRecord
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
								"serviceType":0
							}
						}
					`,
				},
			},
		},
	}

	for _, c := range cases {
		if err := handleEvent(c.Record); err != nil {
			t.Error(err)
		}
	}

}
