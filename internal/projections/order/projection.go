package order

import (
	"fmt"

	es "forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/event"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Projection
type Projection struct {
	db        dynamodbiface.DynamoDBAPI
	tableName *string
}

func NewProjection(db dynamodbiface.DynamoDBAPI, tableName string) *Projection {
	return &Projection{
		db:        db,
		tableName: aws.String(tableName),
	}
}

type Order struct {
	OrderID     string            `json:"orderId"`
	ServiceType model.ServiceType `json:"serviceType"`
}

func (p *Projection) ApplyEvent(e es.Event) error {
	switch d := e.Data.(type) {
	case *event.OrderStartedEvent:
		return p.handleOrderStartedEvent(d)
	case *event.OrderServiceTypeSetEvent:
		return p.handleServiceTypeSetEvent(d)
	default:
		return fmt.Errorf("Unsupported event %s received in ApplyEvent handler of the Order Projection: %+v", d, e)
	}

	return nil
}

func (p *Projection) handleOrderStartedEvent(e *event.OrderStartedEvent) error {
	fmt.Printf("Handling projection for OrderStartedEvent: %+v", e)
	return p.save(&Order{
		OrderID:     e.OrderID,
		ServiceType: e.ServiceType,
	})
}

func (p *Projection) handleServiceTypeSetEvent(e *event.OrderServiceTypeSetEvent) error {
	fmt.Printf("Handling projection for OrderServiceTypeSetEvent: %+v", e)
	return nil
}

func (p *Projection) save(order *Order) error {
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		return err
	}

	_, err = p.db.PutItem(&dynamodb.PutItemInput{
		TableName: p.tableName,
		Item:      av,
	})
	return err
}
