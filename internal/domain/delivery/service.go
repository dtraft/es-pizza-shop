package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/delivery/command"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

var apiURL = "https://jsonplaceholder.cypress.io"

type OrderDelivery struct {
	DeliveryID  int    `json:"id"`
	Description string `json:"description"`
}

type ServiceAPI interface {
	ReceiveDeliveryNotification(int) error
	SubmitOrderForDelivery(*OrderDelivery) (*OrderDelivery, error)
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

func (s *Service) ReceiveDeliveryNotification(deliveryID int) error {
	return s.processCommand(&command.ConfirmDelivery{DeliveryID: deliveryID})
}

func (s *Service) SubmitOrderForDelivery(payload *OrderDelivery) (*OrderDelivery, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/posts", s.url), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	o := &OrderDelivery{}
	if err := json.Unmarshal(respBody, o); err != nil {
		return nil, err
	}

	if err := s.processCommand(&command.RequestDelivery{DeliveryID: o.DeliveryID}); err != nil {
		return nil, err
	}

	log.Printf("Delivery submitted, got ID: %d", o.DeliveryID)

	return o, nil
}

var _ ServiceAPI = (*Service)(nil)
