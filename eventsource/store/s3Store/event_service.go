package s3Store

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type EventService struct {
	bucket            *string
	s3                *s3.S3
	continuationToken *string
}

func New(s3 *s3.S3, bucket string) *EventService {
	return &EventService{
		bucket: aws.String(bucket),
		s3:     s3,
	}
}

func (s *EventService) HasNextPage() bool {
	return s.continuationToken != nil
}

func (s *EventService) GetNextPage() ([]eventsource.Event, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:            s.bucket,
		MaxKeys:           aws.Int64(20),
		Prefix:            aws.String("events"),
		ContinuationToken: s.continuationToken,
	}

	result, err := s.s3.ListObjectsV2(input)
	if err != nil {
		return nil, err
	}

	s.continuationToken = result.NextContinuationToken

	var events []eventsource.Event
	for _, item := range result.Contents {
		log.Printf("Retrieving event with key: %s", *item.Key)
		input := &s3.GetObjectInput{
			Bucket: s.bucket,
			Key:    item.Key,
		}

		output, err := s.s3.GetObject(input)
		if err != nil {
			log.Printf("Error retrieving s3Store object: %s", err)
			return nil, err
		}
		body, err := ioutil.ReadAll(output.Body)
		if err != nil {
			log.Printf("Error reading s3Store object body: %s", err)
			return nil, err
		}
		var rawEvent Event
		if err := json.Unmarshal(body, &rawEvent); err != nil {
			log.Printf("Error unmarshalling json to rawEvent: %s", string(body))
			return nil, err
		}

		eventData, err := eventsource.GetEventOfType(rawEvent.EventType)
		if err != nil {
			if e, ok := err.(*eventsource.UnregisteredEventError); ok {
				log.Printf("Skipping unregistered event of type: %s", e.EventType)
				continue
			}
			return nil, err
		}

		if err := eventData.Load(rawEvent.RawData, rawEvent.EventTypeVersion); err != nil {
			return nil, err
		}

		event := eventsource.Event{
			EventID:           rawEvent.EventID,
			AggregateID:       rawEvent.AggregateID,
			AggregateType:     rawEvent.AggregateType,
			AggregateSequence: rawEvent.AggregateSequence,
			EventTypeVersion:  rawEvent.EventTypeVersion,
			EventType:         rawEvent.EventType,
			Timestamp:         rawEvent.Timestamp,
			Data:              eventData,
		}

		events = append(events, event)
	}

	return events, nil
}

// Event is the S3 represenation of a domain event
type Event struct {
	EventID           string          `json:"eventId"`
	AggregateID       string          `json:"aggregateId"`
	AggregateType     string          `json:"aggregateType"`
	AggregateSequence int             `json:"aggregateSequence"`
	EventType         string          `json:"eventType"`
	EventTypeVersion  int             `json:"eventVersion"`
	Timestamp         time.Time       `json:"eventTimestamp"`
	RawData           json.RawMessage `json:"eventData"`
}
