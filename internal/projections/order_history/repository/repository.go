package repository

import (
	"log"

	. "forge.lmig.com/n1505471/pizza-shop/internal/projections/order_history/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Interface interface {
	Save(order *OrderHistoryRecord) error
}

// The Repository provides a way to persist and retrieve entities from permanent storage
type Repository struct {
	db        dynamodbiface.DynamoDBAPI
	tableName *string
}

func NewRepository(db dynamodbiface.DynamoDBAPI, tableName string) *Repository {
	return &Repository{
		db:        db,
		tableName: aws.String(tableName),
	}
}

/*
 * Write Handlers
 */

func (r *Repository) Save(order *OrderHistoryRecord) error {
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		TableName: r.tableName,
		Item:      av,
	})
	return err
}

/*
 * Query Handlers
 */

// QueryAllOrders retrieves a list of all orders, sorted in random order
func (r *Repository) QueryHistoryForOrderID(orderId string) ([]*OrderHistoryRecord, error) {
	result, err := r.db.Query(&dynamodb.QueryInput{
		TableName: r.tableName,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":orderId": {
				S: aws.String(orderId),
			},
		},
		KeyConditionExpression: aws.String("orderId = :orderId"),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Raw Orders History Records: %+v", result.Items)

	orders := []*OrderHistoryRecord{}
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
