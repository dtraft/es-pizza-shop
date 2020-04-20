package main

import (
	"context"
	"fmt"
	"log"
	"os"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/projections/order"
	"forge.lmig.com/n1505471/pizza-shop/internal/projections/order/repository"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var repo = repository.NewRepository(dynamodb.New(session.New(), aws.NewConfig()), os.Getenv("TABLE_NAME"))
var projection = order.NewProjection(repo)

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

	// Handle projection
	if err := projection.ApplyEvent(event); err != nil {
		return fmt.Errorf("Error handling event with payload: %+v, details: %s", event, err)
	}

	return nil
}
