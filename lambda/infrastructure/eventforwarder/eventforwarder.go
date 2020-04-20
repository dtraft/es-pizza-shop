package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource/store/dynamodb"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

var snsClient snsiface.SNSAPI
var s3Client s3iface.S3API

var eventBus = aws.String(os.Getenv("EVENT_BUS"))
var bucketName = aws.String(os.Getenv("BUCKET_NAME"))

func load() {
	session := session.New()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		snsClient = sns.New(session)
		wg.Done()
	}()

	go func() {
		s3Client = s3.New(session)
		wg.Done()
	}()

	wg.Wait()
}

func main() {
	if snsClient == nil || s3Client == nil {
		load()
	}
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, e DynamoEvent) error {
	for _, r := range e.Records {
		switch r.EventName {
		case "INSERT":
			fallthrough
		case "MODIFY":

			event := es.Event{}
			if err := dynamodbattribute.UnmarshalMap(r.Change.NewImage, &event); err != nil {
				log.Printf("Error decoding event from dynamodb: %s ", err)
				return err
			}

			encoded, err := json.Marshal(event)
			if err != nil {
				log.Printf("Error encoding event to json: %s ", err)
				return err
			}

			wg := &sync.WaitGroup{}
			wg.Add(2)

			// Store event in S3
			go func() {
				key := fmt.Sprintf("events/%s--%s--%s", event.Timestamp.Format("2006-01-02T15:04:05.999Z"), event.EventType, event.EventID)
				i := &s3.PutObjectInput{
					Bucket:      bucketName,
					ContentType: aws.String("application/json"),
					Key:         aws.String(key),
					Body:        aws.ReadSeekCloser(bytes.NewReader(encoded)),
				}

				_, err := s3Client.PutObject(i)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok {
						switch aerr.Code() {
						default:
							log.Printf("Error putting event to s3: %s.\n", aerr.Error())
						}
					} else {
						// Print the error, cast err to awserr.Error to get the Code and
						// Message from an error.
						log.Printf("Error putting event to s3: %s.\n", err.Error())
					}
				}
				wg.Done()
			}()

			// Forward event to event bus
			go func() {

				i := &sns.PublishInput{
					TopicArn: eventBus,
					Message:  aws.String(string(encoded)),
					MessageAttributes: map[string]*sns.MessageAttributeValue{
						"eventType": {
							DataType:    aws.String("String"),
							StringValue: aws.String(event.EventType),
						},
						"eventVersion": {
							DataType:    aws.String("Number"),
							StringValue: aws.String(strconv.Itoa(event.EventTypeVersion)),
						},
						"eventId": {
							DataType:    aws.String("String"),
							StringValue: aws.String(event.EventID),
						},
					},
				}

				log.Printf("Publish Input: %+v", i)

				_, err := snsClient.Publish(i)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok {
						switch aerr.Code() {
						default:
							log.Printf("Error publishing event to sns: %s.\n", aerr.Error())
						}
					} else {
						// Print the error, cast err to awserr.Error to get the Code and
						// Message from an error.
						log.Printf("Error publishing event to sns: %s.\n", err.Error())
					}
				}

				wg.Done()
			}()

			wg.Wait()
		}
	}
	return nil
}

type DynamoEventChange struct {
	NewImage map[string]*dynamodb.AttributeValue `json:"NewImage"`
	// ... more fields if needed: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_streams_GetRecords.html
}

type DynamoEventRecord struct {
	Change    DynamoEventChange `json:"dynamodb"`
	EventName string            `json:"eventName"`
	EventID   string            `json:"eventID"`
	// ... more fields if needed: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_streams_GetRecords.html
}

type DynamoEvent struct {
	Records []DynamoEventRecord `json:"Records"`
}
