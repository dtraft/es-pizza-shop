// generated by jsonenums -type=Status; DO NOT EDIT

package model

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	_StatusNameToValue = map[string]Status{
		"Started":   Started,
		"Submitted": Submitted,
		"Approved":  Approved,
		"Delivered": Delivered,
	}

	_StatusValueToName = map[Status]string{
		Started:   "Started",
		Submitted: "Submitted",
		Approved:  "Approved",
		Delivered: "Delivered",
	}
)

func init() {
	var v Status
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_StatusNameToValue = map[string]Status{
			interface{}(Started).(fmt.Stringer).String():   Started,
			interface{}(Submitted).(fmt.Stringer).String(): Submitted,
			interface{}(Approved).(fmt.Stringer).String():  Approved,
			interface{}(Delivered).(fmt.Stringer).String(): Delivered,
		}
	}
}

// MarshalJSON is generated so Status satisfies json.Marshaler.
func (r Status) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _StatusValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid Status: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so Status satisfies json.Unmarshaler.
func (r *Status) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Status should be a string, got %s", data)
	}
	v, ok := _StatusNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid Status %q", s)
	}
	*r = v
	return nil
}

func (r *Status) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		av.S = aws.String(s.String())
		return nil
	}
	s, ok := _StatusValueToName[*r]
	if !ok {
		return fmt.Errorf("invalid Status: %d", r)
	}
	av.S = aws.String(s)
	return nil
}

func (r *Status) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av.S == nil {
		return nil
	}

	s := aws.StringValue(av.S)
	v, ok := _StatusNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid Status %q", s)
	}
	*r = v
	return nil
}