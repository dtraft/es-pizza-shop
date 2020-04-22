package repository

import (
	"fmt"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	. "forge.lmig.com/n1505471/pizza-shop/internal/projections/order/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-test/deep"
	"io/ioutil"
	"log"
	"os"
	"testing"
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
					":serviceType": {N: aws.String(fmt.Sprintf("%d", model.Pickup))},
				},
				UpdateExpression: aws.String("SET #serviceType = :serviceType"),
			},
		},
	}

	for _, c := range cases {
		mockDb.Expected = c.Expected
		if err := repo.Patch(mockOrderID, c.Updates); err != nil {
			t.Error(err)
		}
	}
}

//
//func TestPatchHelper(t *testing.T) {
//
//	cases := []struct {
//		Updates  interface{}
//		Expected map[string]*dynamodb.AttributeValue
//	}{
//		{
//			Updates: map[string]interface{}{
//				"test": "me",
//				"nested": map[string]interface{}{
//					"test":    "4u",
//					"another": "4us",
//				},
//			},
//			Expected: map[string]*dynamodb.AttributeValue{
//				"test":           {S: aws.String("me")},
//				"nested.test":    {S: aws.String("4u")},
//				"nested.another": {S: aws.String("4us")},
//			},
//		},
//	}
//
//	for _, c := range cases {
//		vals := make(map[string]*dynamodb.AttributeValue)
//
//		if err := patchHelper(c.Updates, "", vals); err != nil {
//			t.Error(err)
//		}
//		if diff := deep.Equal(c.Expected, vals); diff != nil {
//			t.Error(diff)
//		}
//	}
//
//}

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
