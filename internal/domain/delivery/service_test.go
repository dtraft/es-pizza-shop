package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/command"
)

var testDescription = "test description"
var deliveryID = 101

func TestService_ReceiveDeliveryNotification(t *testing.T) {
	cases := []struct {
		Label       string
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should correctly issue the DeliveryConfirmed command",
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.ConfirmDelivery)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.ConfirmDelivery{}, c)
				}
				if cmd.DeliveryID != 101 {
					return fmt.Errorf("Expected `%d` for DeliveryID, got `%d`", 101, cmd.DeliveryID)
				}
				return nil
			},
		},
	}

	for i, c := range cases {
		s := NewService(&mockEventSource{
			check:       c.Check,
			shouldError: c.ShouldError,
		})

		err := s.ReceiveDeliveryNotification(101)
		if c.ShouldError && err == nil {
			t.Errorf("Cases[%d] FAILED: %s, expected an error.", i, c.Label)
			continue
		}

		if !c.ShouldError && err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, err)
		}
	}
}

func TestService_SubmitOrderForApproval(t *testing.T) {
	cases := []struct {
		Label             string
		Payload           *OrderDelivery
		HanderFuncFactory func(t *testing.T, label string, i int) http.HandlerFunc
		Check             Condition
		Expected          *OrderDelivery
		ShouldError       bool
	}{
		{
			Label:             "Should correctly process orders submitted for delivery",
			Payload:           &OrderDelivery{Description: testDescription},
			HanderFuncFactory: makeSuccessRequestDeliveryHandler,
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.RequestDelivery)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.RequestDelivery{}, c)
				}
				if cmd.DeliveryID != 101 {
					return fmt.Errorf("Expected `%d` for DeliveryID, got `%d`", 101, cmd.DeliveryID)
				}
				return nil
			},
			Expected: &OrderDelivery{DeliveryID: deliveryID},
		},
		{
			Label:             "Should return error for failed HTTP calls",
			Payload:           &OrderDelivery{Description: testDescription},
			HanderFuncFactory: makeFailedRequestDeliveryHandler,
			ShouldError:       true,
		},
		{
			Label:             "Should return error when invalid json is returned from HTTP call",
			Payload:           &OrderDelivery{Description: testDescription},
			HanderFuncFactory: makeInvalidRequestDeliveryHandler,
			ShouldError:       true,
		},
		{
			Label:             "Should return error when aggregate issues an error",
			Payload:           &OrderDelivery{Description: testDescription},
			HanderFuncFactory: makeSuccessRequestDeliveryHandler,
			Check: func(c eventsource.Command) error {
				return fmt.Errorf("i am error.")
			},
			ShouldError: true,
		},
	}

	for i, c := range cases {
		ts := httptest.NewServer(c.HanderFuncFactory(t, c.Label, i))

		s := &Service{
			eventSource: &mockEventSource{
				check:       c.Check,
				shouldError: c.ShouldError,
			},
			url: ts.URL,
		}

		result, err := s.SubmitOrderForDelivery(c.Payload)
		if c.ShouldError && err == nil {
			t.Errorf("Cases[%d] FAILED: %s, expected an error.", i, c.Label)
			continue
		}

		if !c.ShouldError && err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, err)
			continue
		}
		if diff := deep.Equal(result, c.Expected); diff != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error: %s", i, c.Label, diff)
		}
		ts.Close()
	}
}

type Condition func(c eventsource.Command) error

type mockEventSource struct {
	eventsource.EventSourceAPI
	check       Condition
	shouldError bool
}

func (m *mockEventSource) ProcessCommand(got eventsource.Command, a eventsource.Aggregate) error {

	if m.shouldError {
		return fmt.Errorf("Expecting an error here.")
	}

	if _, ok := a.(*Aggregate); !ok {
		return fmt.Errorf("Error: Expected %T, got %T", &Aggregate{}, a)
	}

	return m.check(got)
}

func makeSuccessRequestDeliveryHandler(t *testing.T, label string, i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if r.Method != "POST" {
			t.Errorf("Cases[%d] FAILED: %s.  Error: Expected a method of `POST`, got %s", i, label, r.Method)
		}
		if r.URL.EscapedPath() != "/posts" {
			t.Errorf("Cases[%d] FAILED: %s.  Error: Expected a path of `/posts`, got %s", i, label, r.URL.EscapedPath())
		}

		reqJson, err := simplejson.NewFromReader(r.Body)
		if err != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Error reading request JSON: %s", i, label, err)
		}

		description := reqJson.Get("description").MustString()
		if diff := deep.Equal(description, testDescription); diff != nil {
			t.Errorf("Cases[%d] FAILED: %s.  Body `description` field: %s", i, label, diff)
		}

		payload, _ := json.Marshal(map[string]interface{}{
			"id": deliveryID,
		})

		w.Write(payload)
	}
}

func makeFailedRequestDeliveryHandler(t *testing.T, label string, i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(`{"message":"i am error."}`))
	}
}

func makeInvalidRequestDeliveryHandler(t *testing.T, label string, i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte(`{"message":"i am error."`))
	}
}
