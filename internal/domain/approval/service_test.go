package approval

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitly/go-simplejson"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/command"
	"github.com/go-test/deep"
)

var testDescription = "test description"
var approvalID = 101

func TestService_ReceiveApproval(t *testing.T) {
	cases := []struct {
		Label       string
		Check       Condition
		ShouldError bool
	}{
		{
			Label: "Should correctly issue the ReceiveApproval command",
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.ReceiveApproval)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.ReceiveApproval{}, c)
				}
				if cmd.ApprovalID != 101 {
					return fmt.Errorf("Expected `%d` for ApprovalID, got `%d`", 101, cmd.ApprovalID)
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

		err := s.ReceiveApproval(approvalID)
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
		Payload           *OrderApproval
		HanderFuncFactory func(t *testing.T, label string, i int) http.HandlerFunc
		Check             Condition
		Expected          *OrderApproval
		ShouldError       bool
	}{
		{
			Label:             "Should correctly process orders submitted for approval",
			Payload:           &OrderApproval{Description: testDescription},
			HanderFuncFactory: makeSuccessRequestApprovalHandler,
			Check: func(c eventsource.Command) error {
				cmd, ok := c.(*command.RequestApproval)
				if !ok {
					return fmt.Errorf("Expected %T, got %T", &command.RequestApproval{}, c)
				}
				if cmd.ApprovalID != 101 {
					return fmt.Errorf("Expected `%d` for ApprovalID, got `%d`", 101, cmd.ApprovalID)
				}
				return nil
			},
			Expected: &OrderApproval{ApprovalID: approvalID},
		},
		{
			Label:             "Should return error for failed HTTP calls",
			Payload:           &OrderApproval{Description: testDescription},
			HanderFuncFactory: makeFailedRequestApprovalHandler,
			ShouldError:       true,
		},
		{
			Label:             "Should return error when invalid json is returned from HTTP call",
			Payload:           &OrderApproval{Description: testDescription},
			HanderFuncFactory: makeInvalidRequestApprovalHandler,
			ShouldError:       true,
		},
		{
			Label:             "Should return error when aggregate issues an error",
			Payload:           &OrderApproval{Description: testDescription},
			HanderFuncFactory: makeSuccessRequestApprovalHandler,
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

		result, err := s.SubmitOrderForApproval(c.Payload)
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

func makeSuccessRequestApprovalHandler(t *testing.T, label string, i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if r.Method != "POST" {
			t.Errorf("Cases[%d] FAILED: %s.  Error: Expected a method of `POST`, got %s", i, label, r.Method)
		}
		if r.URL.EscapedPath() != "/todos" {
			t.Errorf("Cases[%d] FAILED: %s.  Error: Expected a path of `/todos`, got %s", i, label, r.URL.EscapedPath())
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
			"id": approvalID,
		})

		w.Write(payload)
	}
}

func makeFailedRequestApprovalHandler(t *testing.T, label string, i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(`{"message":"i am error."}`))
	}
}

func makeInvalidRequestApprovalHandler(t *testing.T, label string, i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte(`{"message":"i am error."`))
	}
}
