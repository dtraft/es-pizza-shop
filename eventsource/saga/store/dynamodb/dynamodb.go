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
	CompositeKey string `dynamodbav:"compositeKey"`
	SagaId       string `dynamodbav:"sagaId"`
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

func (s *SagaStore) Load(association *saga.SagaAssociation, sagaType string) (*saga.Wrapper, error) {
	sagaId, err := s.retrieveSagaID(association, sagaType)
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

	return &saga.Wrapper{
		ID:      out.ID,
		Version: out.Version,
		Data:    encoded,
	}, nil
}

func (s *SagaStore) AddAssociationID(association *saga.SagaAssociation, wrapper *saga.Wrapper) error {

	compositeKey := fmt.Sprintf("%s#%s#%s", association.ID, association.AssociationType, wrapper.Type)
	av, err := dynamodbattribute.MarshalMap(&sagaAssociation{
		CompositeKey: compositeKey,
		SagaId:       wrapper.ID,
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

func (s *SagaStore) Save(wrapper *saga.Wrapper) error {
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

func (s *SagaStore) retrieveSagaID(association *saga.SagaAssociation, sagaType string) (string, error) {
	compositeKey := fmt.Sprintf("%s#%s#%s", association.ID, association.AssociationType, sagaType)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"compositeKey": {
				S: aws.String(compositeKey),
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
					AssociationID: association.ID,
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

var _ saga.Storer = (*SagaStore)(nil)
