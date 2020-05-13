package approval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/command"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

var apiURL = "https://jsonplaceholder.cypress.io"

type OrderApproval struct {
	ApprovalID  int    `json:"id"`
	Description string `json:"description"`
}

type ServiceAPI interface {
	ReceiveApproval(int) error
	SubmitOrderForApproval(*OrderApproval) (*OrderApproval, error)
}

type Service struct {
	eventSource eventsource.EventSourceAPI
	url         string
}

func (s *Service) processCommand(c eventsource.Command) error {
	return s.eventSource.ProcessCommand(c, &Aggregate{})
}

func NewService(eventSource eventsource.EventSourceAPI) *Service {
	return &Service{
		eventSource: eventSource,
		url:         apiURL,
	}
}

func (s *Service) ReceiveApproval(approvalID int) error {
	return s.processCommand(&command.ReceiveApproval{ApprovalID: approvalID})
}

func (s *Service) SubmitOrderForApproval(payload *OrderApproval) (*OrderApproval, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/todos", s.url), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("Approval service returned %d status code, details: %s", resp.StatusCode, respBody)
	}

	o := &OrderApproval{}
	if err := json.Unmarshal(respBody, o); err != nil {
		return nil, err
	}

	if err := s.processCommand(&command.RequestApproval{ApprovalID: o.ApprovalID}); err != nil {
		return nil, err
	}
	log.Printf("Approval requested with payload: %+v, got tracking ID: %d", payload, o.ApprovalID)

	return o, nil
}

var _ ServiceAPI = (*Service)(nil)
