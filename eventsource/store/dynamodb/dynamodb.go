package dynamodb

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type EventStore struct {
	svc       *dynamodb.DynamoDB
	tableName *string
}

func New(svc *dynamodb.DynamoDB, t string) *EventStore {
	return &EventStore{
		svc:       svc,
		tableName: aws.String(t),
	}
}

func (e *EventStore) SaveEvent(event eventsource.Event) error {
	av, err := dynamodbattribute.MarshalMap(event)
	if err != nil {
		return err
	}

	err = e.save(av)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return &eventsource.AggregateLockError{
					ID:       event.AggregateID,
					Sequence: event.AggregateSequence,
				}
			}
		}
	}

	return err
}

func (e *EventStore) EventsForAggregate(aggregateID string) ([]eventsource.Event, error) {
	var events []eventsource.Event
	av, err := dynamodbattribute.Marshal(aggregateID)
	if err != nil {
		return events, err
	}
	results, err := e.query("aggregateId = :aggregateId", map[string]*dynamodb.AttributeValue{
		":aggregateId": av,
	})
	return unmarshalEventsFromDB(results)
}

func (e *EventStore) save(item map[string]*dynamodb.AttributeValue) error {
	_, err := e.svc.PutItem(&dynamodb.PutItemInput{
		TableName:           e.tableName,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(aggregateSequence)"),
	})

	return err
}

func (e *EventStore) query(query string, attributeValues map[string]*dynamodb.AttributeValue) ([]map[string]*dynamodb.AttributeValue, error) {
	result, err := e.svc.Query(&dynamodb.QueryInput{
		KeyConditionExpression:    aws.String(query),
		ExpressionAttributeValues: attributeValues,
		TableName:                 e.tableName,
	})
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// Event is the DynamoDB represenation of a domain event
type Event struct {
	EventID           string                 `json:"eventId"`
	AggregateID       string                 `json:"aggregateId"`
	AggregateType     string                 `json:"aggregateType"`
	AggregateSequence int                    `json:"aggregateSequence"`
	EventType         string                 `json:"eventType"`
	EventTypeVersion  int                    `json:"eventVersion"`
	Timestamp         time.Time              `json:"eventTimestamp"`
	RawData           map[string]interface{} `json:"eventData"`
}

func (e *Event) toDomainEvent() (eventsource.Event, error) {
	eventData, err := eventsource.GetEventOfType(e.EventType)
	if err != nil {
		return eventsource.Event{}, err
	}

	r, err := json.Marshal(e.RawData)
	if err != nil {
		return eventsource.Event{}, err
	}
	if err := eventData.Load(r, e.EventTypeVersion); err != nil {
		return eventsource.Event{}, err
	}
	event := eventsource.Event{
		EventID:           e.EventID,
		AggregateID:       e.AggregateID,
		AggregateType:     e.AggregateType,
		AggregateSequence: e.AggregateSequence,
		EventTypeVersion:  e.EventTypeVersion,
		EventType:         e.EventType,
		Timestamp:         e.Timestamp,
		Data:              eventData,
	}

	return event, nil
}

func unmarshalEventsFromDB(results []map[string]*dynamodb.AttributeValue) ([]eventsource.Event, error) {
	var events []eventsource.Event
	dbEvents := []Event{}
	err := dynamodbattribute.UnmarshalListOfMaps(results, &dbEvents)
	if err != nil {
		return events, err
	}
	events = make([]eventsource.Event, len(dbEvents))
	for i, dbEvent := range dbEvents {
		event, err := dbEvent.toDomainEvent()
		if err != nil {
			return events, err
		}
		events[i] = event
	}

	return events, nil
}
