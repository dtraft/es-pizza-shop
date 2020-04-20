package repository

import (
	"fmt"
	"log"
	"strings"

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

func (r *Repository) Patch(orderID string, updates map[string]interface{}) error {

	vals := make(map[string]*dynamodb.AttributeValue)

	if err := patchHelper(updates, "", vals); err != nil {
		return err
	}

	names := make(map[string]*string)
	values := make(map[string]*dynamodb.AttributeValue)
	var expressions []string
	for name, value := range vals {
		keys := strings.Split(name, ".")
		valueKey := ":" + strings.Join(keys, "")
		for i, key := range keys {
			keys[i] = "#" + key
			names[keys[i]] = aws.String(key)
		}
		nameKey := strings.Join(keys, ".")
		values[valueKey] = value

		expressions = append(expressions, fmt.Sprintf("%s = %s", nameKey, valueKey))
	}

	exp := "SET " + strings.Join(expressions, ", ")

	i := &dynamodb.UpdateItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"orderId": {S: aws.String(orderID)},
		},
		ExpressionAttributeValues: values,
		ExpressionAttributeNames:  names,
		UpdateExpression:          aws.String(exp),
	}

	log.Printf("Update item request: %+v", i)

	if _, err := r.db.UpdateItem(i); err != nil {
		return err
	}
	return nil
}

func patchHelper(update interface{}, path string, vals map[string]*dynamodb.AttributeValue) error {
	switch u := update.(type) {
	case map[string]interface{}:
		for key, value := range u {
			nestedPath := path
			if nestedPath != "" {
				nestedPath += "."
			}
			nestedPath += key
			if err := patchHelper(value, nestedPath, vals); err != nil {
				return err
			}
		}
	default:
		a, err := dynamodbattribute.Marshal(update)
		if err != nil {
			return err
		}
		vals[path] = a
	}
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
