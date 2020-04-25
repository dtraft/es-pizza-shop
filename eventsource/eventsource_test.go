package eventsource

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
)

func init() {
	RegisterEventType(&TestData{})
}

func TestEvent_Load(t *testing.T) {
	cases := []struct {
		Label       string
		Data        string
		Expected    interface{}
		ShouldError bool
	}{
		{
			Label: "Should properly unmarhsall event data",
			Data: `
				{
					"eventId":"6c4539e3-ae1b-44f0-bfc2-4d7531893136",
					"aggregateId":"84de2628-ac3b-4fcf-b2a1-05cf5b1b5743",
					"aggregateType":"TestAggregate",
					"eventVersion":1,
					"eventType":"TestData",
					"eventTimestamp":"2020-04-19T19:45:11.475995951Z",
					"eventData":{
						"testId": "123abc",
						"name": "Me!"
					}
				}
			`,
			Expected: &TestData{
				TestID: "123abc",
				Name:   "Me!",
			},
		},
	}

	for i, c := range cases {
		event := &Event{
			EventType: "TestData",
		}
		err := event.Load([]byte(c.Data))

		if c.ShouldError && err == nil {
			t.Errorf("Cases[%d] FAILED: %s.  Expected an error.", i, c.Label)
		}

		if err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, err)
		}

		if diff := deep.Equal(event.Data, c.Expected); diff != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, diff)
		}
	}
}

/*
 * Set up
 */

type TestAggregate struct {
	Aggregate
	TestID string
}

func (m *TestAggregate) Type() string {
	return "TestAggregate"
}

type TestData struct {
	TestID string `json:"testId"`
	Name   string `json:"name"`
}

func (e *TestData) Version() int {
	return 1
}

func (e *TestData) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}
