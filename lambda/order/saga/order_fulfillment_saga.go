package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"forge.lmig.com/n1505471/pizza-shop/internal/saga/orderfulfillment"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"
	ddbStore "forge.lmig.com/n1505471/pizza-shop/eventsource/saga/store/dynamodb"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var manager *saga.SagaManager

func init() {
	db := dynamodb.New(session.New(), aws.NewConfig())
	store := ddbStore.New(db, os.Getenv("ASSOCIATIONS_TABLE_NAME"), os.Getenv("SAGA_TABLE_NAME"))
	manager = saga.NewManager(store)
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, e events.SNSEvent) error {

	for _, r := range e.Records {
		if err := handleEvent(r); err != nil {
			log.Println(err)
			continue
		}
	}

	return nil
}

func handleEvent(r events.SNSEventRecord) error {
	eventTypeAttribute := r.SNS.MessageAttributes["eventType"].(map[string]interface{})
	event := es.Event{
		EventType: eventTypeAttribute["Value"].(string),
	}
	if err := event.Load([]byte(r.SNS.Message)); err != nil {
		return fmt.Errorf("Error unmarhalling json: %s", err)
	}

	// Handle saga
	if err := manager.ProcessEvent(event, &orderfulfillment.OrderFulfillmentSaga{}); err != nil {
		return fmt.Errorf("Error handling event with payload: %+v, details: %s", event, err)
	}

	return nil
}
