package approval

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/approval/command"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
)

var url = "https://jsonplaceholder.cypress.io/todos"

type OrderApproval struct {
	ApprovalID  int               `json:"id"`
	ServiceType model.ServiceType `json:"serviceType"`
	Description string            `json:"description"`
}

type ServiceAPI interface {
	ReceiveApproval(int) error
	SubmitOrderForApproval(*OrderApproval) (*OrderApproval, error)
}

type Service struct {
	eventSource eventsource.EventSourceAPI
}

func (s *Service) processCommand(c eventsource.Command) error {
	return s.eventSource.ProcessCommand(c, &Aggregate{})
}

func NewService(eventSource eventsource.EventSourceAPI) *Service {
	return &Service{
		eventSource: eventSource,
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
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var o *OrderApproval
	if err := json.Unmarshal(respBody, o); err != nil {
		return nil, err
	}

	if err := s.processCommand(&command.RequestApproval{ApprovalID: o.ApprovalID}); err != nil {
		return nil, err
	}

	return o, nil
}

var _ ServiceAPI = (*Service)(nil)
