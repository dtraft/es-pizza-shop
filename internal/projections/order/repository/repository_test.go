package repository

import (
	"fmt"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-test/deep"
	"testing"
)

type mockDynamoDb struct {
	dynamodbiface.DynamoDBAPI
	Expected interface{}
}

var mockTable = "test"
var mockOrderID = "testOrderId"
var mockDb = &mockDynamoDb{}
var repo = NewRepository(mockDb, mockTable)

func (m mockDynamoDb) UpdateItem(in *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	// Check if this matches the expected
	if diff := deep.Equal(m.Expected, in); diff != nil {
		return nil, fmt.Errorf("%s", diff)
	}
	// Only need to return mocked response output
	return &dynamodb.UpdateItemOutput{}, nil
}

func TestRepository_Patch(t *testing.T) {
	cases := []struct {
		Updates  map[string]interface{}
		Expected *dynamodb.UpdateItemInput
	}{
		{
			Updates: map[string]interface{}{
				"test": "me",
				"nested": map[string]interface{}{
					"test":    "4u",
					"another": "4us",
				},
			},
			Expected: &dynamodb.UpdateItemInput{
				TableName: aws.String(mockTable),
				Key: map[string]*dynamodb.AttributeValue{
					"orderId": {S: aws.String(mockOrderID)},
				},
				ExpressionAttributeNames: map[string]*string{
					"#nested":  aws.String("nested"),
					"#another": aws.String("another"),
					"#test":    aws.String("test"),
				},
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":nestedanother": {S: aws.String("4us")},
					":test":          {S: aws.String("me")},
					":nestedtest":    {S: aws.String("4u")},
				},
				UpdateExpression: aws.String("SET #test = :test, #nested.#test = :nestedtest, #nested.#another = :nestedanother"),
			},
		},
		{
			Updates: map[string]interface{}{
				"serviceType": model.Delivery,
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
					":serviceType": {N: aws.String(fmt.Sprintf("%d", model.Delivery))},
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

func TestPatchHelper(t *testing.T) {

	cases := []struct {
		Updates  interface{}
		Expected map[string]*dynamodb.AttributeValue
	}{
		{
			Updates: map[string]interface{}{
				"test": "me",
				"nested": map[string]interface{}{
					"test":    "4u",
					"another": "4us",
				},
			},
			Expected: map[string]*dynamodb.AttributeValue{
				"test":           {S: aws.String("me")},
				"nested.test":    {S: aws.String("4u")},
				"nested.another": {S: aws.String("4us")},
			},
		},
	}

	for _, c := range cases {
		vals := make(map[string]*dynamodb.AttributeValue)

		if err := patchHelper(c.Updates, "", vals); err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(c.Expected, vals); diff != nil {
			t.Error(diff)
		}
	}

}
