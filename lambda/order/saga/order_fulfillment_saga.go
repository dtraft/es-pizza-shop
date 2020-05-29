package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval"

	"forge.lmig.com/n1505471/pizza-shop/internal/saga/orderfulfillment"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"
	ddbSagaStore "forge.lmig.com/n1505471/pizza-shop/eventsource/saga/store/dynamodb"
	ddbEventStore "forge.lmig.com/n1505471/pizza-shop/eventsource/store/dynamodb"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var manager *saga.SagaManager
var deliverySvc delivery.ServiceAPI
var approvalSvc approval.ServiceAPI
var orderSvc order.ServiceAPI
var eventsource es.EventSourceAPI

func init() {
	db := dynamodb.New(session.New(), aws.NewConfig())
	store := ddbSagaStore.New(db, os.Getenv("ASSOCIATIONS_TABLE_NAME"), os.Getenv("SAGA_TABLE_NAME"))
	eventStore := ddbEventStore.New(db, os.Getenv("EVENT_TABLE_NAME"))
	eventsource = es.New(eventStore)
	manager = saga.NewManager(store)
	deliverySvc = delivery.NewService(eventsource)
	approvalSvc = approval.NewService(eventsource)
	orderSvc = order.NewService(eventsource)
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
	orderFulfillmentSaga := orderfulfillment.New(orderSvc, deliverySvc, approvalSvc)
	if err := manager.ProcessEvent(event, orderFulfillmentSaga); err != nil {
		return fmt.Errorf("Error handling event with payload: %+v, details: %s", event, err)
	}

	return nil
}
