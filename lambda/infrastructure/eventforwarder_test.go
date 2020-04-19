package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource/store/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/google/go-cmp/cmp"
)

type mockedPutRecord struct {
	kinesisiface.KinesisAPI
	Expected []kinesis.PutRecordInput
}

func (m mockedPutRecord) PutRecord(in *kinesis.PutRecordInput) (*kinesis.PutRecordOutput, error) {
	expected := m.Expected[counter]

	if cmp.Equal(in, m.Expected) {
		return nil, fmt.Errorf("In PutRecord, expected %+v, got %+v", expected, in)
	}

	o := &kinesis.PutRecordOutput{
		SequenceNumber: aws.String(strconv.Itoa(counter)),
	}

	counter++
	return o, nil
}

var counter int

func TestHandleRequest(t *testing.T) {

	records := []es.Event{
		{
			EventID:          "eventId",
			AggregateID:      "aggregateId",
			AggregateType:    "aggregateType",
			EventTypeVersion: 1,
			EventType:        "TestType",
			Timestamp:        time.Now(),
			RawData: map[string]interface{}{
				"test": "test",
			},
		},
	}

	var encoded []*dynamodb.AttributeValue
	var raw [][]byte
	for _, r := range records {
		av, err := dynamodbattribute.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		encoded = append(encoded, av)

		j, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		raw = append(raw, j)
	}

	cases := []struct {
		Event    DynamoEvent
		Expected []kinesis.PutRecordInput
	}{
		{
			Event: DynamoEvent{
				Records: []DynamoEventRecord{
					{
						EventName: "INSERT",
						Change: DynamoEventChange{
							NewImage: encoded[0],
						},
					},
				},
			},
			Expected: []kinesis.PutRecordInput{
				{
					PartitionKey: aws.String("aggregateId"),
					StreamName:   aws.String("Events"),
					Data:         raw[0],
				},
			},
		},
	}

	for _, c := range cases {
		svc = &mockedPutRecord{
			Expected: c.Expected,
		}

		err := HandleRequest(nil, c.Event)

		if err != nil {
			t.Error(err)
		}
	}

}
