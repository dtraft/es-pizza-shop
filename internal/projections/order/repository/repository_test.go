package repository

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	. "forge.lmig.com/n1505471/pizza-shop/internal/projections/order/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-test/deep"
)

var mockTable = "test"
var mockOrderID = "testOrderId"
var mockDb = &mockDynamoDb{}
var repo = NewRepository(mockDb, mockTable)

// SETUP
func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestRepository_Patch(t *testing.T) {
	cases := []struct {
		Updates  *Order
		Expected *dynamodb.UpdateItemInput
	}{
		{
			Updates: &Order{
				ServiceType: model.Pickup,
			},
			Expected: &dynamodb.UpdateItemInput{
				TableName: aws.String(mockTable),
				Key: map[string]*dynamodb.AttributeValue{
					"orderId": {S: aws.String(mockOrderID)},
				},
				ExpressionAttributeNames: map[string]*string{
					"#serviceType": aws.String("serviceType"),
				},
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":serviceType": {S: aws.String("Pickup")},
				},
				UpdateExpression: aws.String("SET #serviceType = :serviceType"),
			},
		},
		{
			Updates: &Order{
				Description: "I'm a test!",
			},
			Expected: &dynamodb.UpdateItemInput{
				TableName: aws.String(mockTable),
				Key: map[string]*dynamodb.AttributeValue{
					"orderId": {S: aws.String(mockOrderID)},
				},
				ExpressionAttributeNames: map[string]*string{
					"#description": aws.String("description"),
				},
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":description": {S: aws.String("I'm a test!")},
				},
				UpdateExpression: aws.String("SET #description = :description"),
			},
		},
		{
			Updates: &Order{
				Status: model.Submitted,
			},
			Expected: &dynamodb.UpdateItemInput{
				TableName: aws.String(mockTable),
				Key: map[string]*dynamodb.AttributeValue{
					"orderId": {S: aws.String(mockOrderID)},
				},
				ExpressionAttributeNames: map[string]*string{
					"#status": aws.String("status"),
				},
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":status": {S: aws.String("Submitted")},
				},
				UpdateExpression: aws.String("SET #status = :status"),
			},
		},
	}

	for i, c := range cases {
		mockDb.Expected = c.Expected
		if err := repo.Patch(mockOrderID, c.Updates); err != nil {
			t.Errorf("Cases[%d]: %s", i, err)
		}
	}
}

type mockDynamoDb struct {
	dynamodbiface.DynamoDBAPI
	Expected interface{}
}

func (m mockDynamoDb) UpdateItem(in *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	// Check if this matches the expected
	if diff := deep.Equal(m.Expected, in); diff != nil {
		return nil, fmt.Errorf("%s", diff)
	}
	// Only need to return mocked response output
	return &dynamodb.UpdateItemOutput{}, nil
}
