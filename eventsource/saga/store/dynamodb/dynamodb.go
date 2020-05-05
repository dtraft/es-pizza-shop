package dynamodb

import (
	"encoding/json"
	"fmt"

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
	AssociationIdSagaTypeCompositeKey string `dynamodbav:"associationIdSagaTypeCompositeKey"`
	SagaId                            string `dynamodbav:"sagaId"`
}

type sagaDto struct {
	ID      string
	Version int
	Data    interface{}
}

func New(svc *dynamodb.DynamoDB, associationsTable string, sagaTable string) *SagaStore {
	return &SagaStore{
		svc:               svc,
		associationsTable: aws.String(associationsTable),
		sagaTable:         aws.String(sagaTable),
	}
}

func (s *SagaStore) Load(associationID string, sagaType string) (*saga.RawSagaWrapper, error) {
	sagaId, err := s.retrieveSagaID(associationID, sagaType)
	if err != nil {
		return nil, err
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
				return nil, &saga.SagaNotFoundError{
					SagaID: sagaId,
				}
			}
		}
		return nil, err
	}

	out := &sagaDto{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, out); err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}

	return &saga.RawSagaWrapper{
		ID:      out.ID,
		Version: out.Version,
		Data:    encoded,
	}, nil
}

func (s *SagaStore) AddAssociationID(associationID string, wrapper *saga.RawSagaWrapper) error {

	compositeKey := fmt.Sprintf("%s#%s", associationID, wrapper.Type)
	av, err := dynamodbattribute.MarshalMap(&sagaAssociation{
		AssociationIdSagaTypeCompositeKey: compositeKey,
		SagaId:                            wrapper.ID,
	})
	if err != nil {
		return err
	}
	_, err = s.svc.PutItem(&dynamodb.PutItemInput{
		TableName: s.associationsTable,
		Item:      av,
	})
	if err != nil {
		return err
	}
}

func (s *SagaStore) Save(wrapper *saga.RawSagaWrapper) error {
	var out interface{}
	if err := json.Unmarshal(wrapper.Data, &out); err != nil {
		return err
	}

	av, err := dynamodbattribute.MarshalMap(&sagaDto{
		ID:      wrapper.ID,
		Version: wrapper.Version,
		Data:    out,
	})
	if err != nil {
		return err
	}
	_, err = s.svc.PutItem(&dynamodb.PutItemInput{
		TableName: s.associationsTable,
		Item:      av,
	})
	if err != nil {
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

var _ saga.SagaStorer = (*SagaStore)(nil)
