package repository

import (
	"log"

	. "forge.lmig.com/n1505471/pizza-shop/internal/projections/order/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Interface interface {
	Save(order *Order) error
	Patch(orderID string, updates map[string]interface{}) error
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

func (r *Repository) Save(order *Order) error {
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

func (r *Repository) Patch(orderId string, updates map[string]interface{}) error {

	// av, err := dynamodbattribute.MarshalMap(order)
	// if err != nil {
	// 	return err
	// }

	// _, err = p.db.PutItem(&dynamodb.PutItemInput{
	// 	TableName: p.tableName,
	// 	Item:      av,
	// })
	// return err
	return nil
}

/*
 * Query Handlers
 */

// QueryAllOrders retrieves a list of all orders, sorted in random order
func (r *Repository) QueryAllOrders() ([]Order, error) {

	result, err := r.db.Scan(&dynamodb.ScanInput{
		TableName: r.tableName,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Raw Orders: %+v", result.Items)

	orders := []Order{}
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
