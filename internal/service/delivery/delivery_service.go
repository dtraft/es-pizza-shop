package delivery

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var url = "https://jsonplaceholder.cypress.io/posts"

type OrderDelivery struct {
	DeliveryID  int    `json:"id"`
	Description string `json:"description"`
	Address     string `json:"address"`
}

func SubmitOrderForApproval(payload *OrderDelivery) (*OrderDelivery, error) {
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

	var o *OrderDelivery
	if err := json.Unmarshal(respBody, o); err != nil {
		return nil, err
	}

	return o, nil
}
