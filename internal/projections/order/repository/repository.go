package query

import (
	"log"

	. "forge.lmig.com/n1505471/pizza-shop/internal/projections/order/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Projection
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
