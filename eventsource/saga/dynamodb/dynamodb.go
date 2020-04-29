package dynamodb

import (
	"forge.lmig.com/n1505471/pizza-shop/eventsource/saga"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type SagaStore struct {
	svc               *dynamodb.DynamoDB
	associationsTable *string
	sagaTable         *string
}

type sagaAssociation struct {
	AssociationId string `dynamodbav:"associationId"`
	SagaType      string `dynamodbav:"sagaType"`
	SagaId        string `dynamodbav:"sagaId"`
}

func New(svc *dynamodb.DynamoDB, associationsTable string, sagaTable string) *SagaStore {
	return &SagaStore{
		svc:               svc,
		associationsTable: aws.String(associationsTable),
		sagaTable:         aws.String(sagaTable),
	}
}

func (s *SagaStore) Load(associationID string, sagaType string, in interface{}) error {
	sagaId, err := s.retrieveSagaID(associationID, sagaType)
	if err != nil {
		return err
	}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"sagaId": {
				S: aws.String(sagaId),
			},
		},
		TableName: s.sagaTable,
	}

	result, err := s.svc.GetItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return &saga.SagaNotFoundError{
					SagaID: sagaId,
				}
			}
		}
		return err
	}

	if err := dynamodbattribute.UnmarshalMap(result.Item, in); err != nil {
		return err
	}

	return nil
}

func (s *SagaStore) retrieveSagaID(associationID string, sagaType string) (string, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"associationId": {
				S: aws.String(associationID),
			},
			"sagaType": {
				S: aws.String(sagaType),
			},
		},
		TableName: s.associationsTable,
	}

	result, err := s.svc.GetItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return "", &saga.SagaAssociationNotFoundError{
					AssociationID: associationID,
					SagaType:      sagaType,
				}
			}
		}
		return "", err
	}

	a := &sagaAssociation{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, a); err != nil {
		return "", err
	}

	return a.SagaId, nil
}
