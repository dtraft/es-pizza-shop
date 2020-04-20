package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/go-test/deep"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource/store/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// SETUP
func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestHandleRequest(t *testing.T) {
	// Override for testing
	bucketName = aws.String("testBucket")
	eventBus = aws.String("eventBus")

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

	var encoded []map[string]*dynamodb.AttributeValue
	var raw [][]byte
	for _, r := range records {
		av, err := dynamodbattribute.MarshalMap(r)
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
		Event       DynamoEvent
		ExpectedS3  []*s3.PutObjectInput
		ExpectedSNS []*sns.PublishInput
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
			ExpectedS3: []*s3.PutObjectInput{
				{
					Bucket:      bucketName,
					ContentType: aws.String("application/json"),
					Key:         aws.String(fmt.Sprintf("events/%s--%s--%s", records[0].Timestamp.Format("2006-01-02T15:04:05.999Z"), records[0].EventType, records[0].EventID)),
					Body:        aws.ReadSeekCloser(bytes.NewReader(raw[0])),
				},
			},
			ExpectedSNS: []*sns.PublishInput{
				{
					TopicArn: eventBus,
					Message:  aws.String(string(raw[0])),
					MessageAttributes: map[string]*sns.MessageAttributeValue{
						"eventType": {
							DataType:    aws.String("String"),
							StringValue: aws.String(records[0].EventType),
						},
						"eventVersion": {
							DataType:    aws.String("Number"),
							StringValue: aws.String(strconv.Itoa(records[0].EventTypeVersion)),
						},
						"eventId": {
							DataType:    aws.String("String"),
							StringValue: aws.String(records[0].EventID),
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		s3Client = &mockedS3Client{
			Expected: c.ExpectedS3,
			t:        t,
		}

		snsClient = &mockedSNSClient{
			Expected: c.ExpectedSNS,
			t:        t,
		}

		err := HandleRequest(nil, c.Event)

		if err != nil {
			t.Error(err)
		}
	}

}

type mockedS3Client struct {
	s3iface.S3API
	Expected []*s3.PutObjectInput
	index    int
	t        *testing.T
}

func (m mockedS3Client) PutObject(in *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	expected := m.Expected[m.index]
	m.index++

	if expected == nil {
		// When expected is nil, that means we should return an error
		return &s3.PutObjectOutput{}, fmt.Errorf("I am error")
	}
	if diff := deep.Equal(expected, in); diff != nil {
		m.t.Error(diff)
	}

	return &s3.PutObjectOutput{}, nil
}

type mockedSNSClient struct {
	snsiface.SNSAPI
	Expected []*sns.PublishInput
	index    int
	t        *testing.T
}

func (m mockedSNSClient) Publish(in *sns.PublishInput) (*sns.PublishOutput, error) {
	expected := m.Expected[m.index]
	m.index++

	if expected == nil {
		// When expected is nil, that means we should return an error
		return &sns.PublishOutput{}, fmt.Errorf("I am error")
	}
	if diff := deep.Equal(expected, in); diff != nil {
		m.t.Error(diff)
	}

	return &sns.PublishOutput{}, nil
}
