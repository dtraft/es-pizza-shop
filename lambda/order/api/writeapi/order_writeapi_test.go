package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

	"github.com/julienschmidt/httprouter"
)

type Condition func(rr *httptest.ResponseRecorder) error

func TestStartOrder(t *testing.T) {
	cases := []struct {
		svc       order.ServiceAPI
		body      string
		condition Condition
	}{
		{
			svc: &mockOrderService{},
			body: `
					"serviceType": 0,
					"description": "Here is a description."
				}
			`,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusBadRequest, rr); err != nil {
					return err
				}
				return nil
			},
		},
		{
			svc: &mockOrderService{},
			body: `
				{
					"serviceType": 0,
					"description": "Here is a description."
				}
			`,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusBadRequest, rr); err != nil {
					return err
				}
				return nil
			},
		},
		{
			svc: &mockOrderService{
				err: fmt.Errorf("Something horrible happened."),
			},
			body: `
				{
					"serviceType": 1,
					"description": "Here is a description."
				}
			`,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusBadRequest, rr); err != nil {
					return err
				}
				return nil
			},
		},
		{
			svc: &mockOrderService{},
			body: `
				{
					"serviceType": 1,
					"description": "Here is a description."
				}
			`,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusOK, rr); err != nil {
					return err
				}

				if err := checkHeader("content-type", "application/json", rr); err != nil {
					return err
				}

				expected := &response{
					OK: true,
					Result: &orderResource{
						OrderID:     "orderId",
						ServiceType: 1,
						Description: "Here is a description.",
					},
				}
				if err := checkResponseBody(expected, &response{Result: &orderResource{}}, rr); err != nil {
					return err
				}
				return nil
			},
		},
	}

	for i, c := range cases {
		// Routing Set up
		con := &Controller{
			orderSvc: c.svc,
		}
		router := httprouter.New()
		con.registerRoutes(router)

		// Request Set Up
		req, _ := http.NewRequest("POST", "/orders", strings.NewReader(c.body))
		req.Header.Set("Content-Type", "application/json")

		// Run
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Evaluate
		if err := c.condition(rr); err != nil {
			t.Errorf("Case[%d]: %s", i, err)
		}
	}
}

func TestToggleOrder(t *testing.T) {

	validReq, _ := http.NewRequest("POST", "/orders/orderId/toggle", &bytes.Reader{})

	cases := []struct {
		svc       order.ServiceAPI
		request   *http.Request
		condition Condition
	}{
		{
			svc:     &mockOrderService{},
			request: validReq,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusOK, rr); err != nil {
					return err
				}
				return nil
			},
		},
		{
			svc: &mockOrderService{
				err: fmt.Errorf("Something horrible happened."),
			},
			request: validReq,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusBadRequest, rr); err != nil {
					return err
				}
				return nil
			},
		},
	}

	for i, c := range cases {
		// Routing Set up
		con := &Controller{
			orderSvc: c.svc,
		}
		router := httprouter.New()
		con.registerRoutes(router)

		// Request Set Up
		c.request.Header.Set("Content-Type", "application/json")

		// Run
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, c.request)

		// Evaluate
		if err := c.condition(rr); err != nil {
			t.Errorf("Case[%d]: %s", i, err)
		}
	}
}

type mockOrderService struct {
	order.ServiceAPI
	err error
}

func (m *mockOrderService) StartOrder(order *model.Order) (string, error) {
	return "orderId", m.err
}

func (m *mockOrderService) ToggleOrderServiceType(orderID string) error {
	return m.err
}

/*
 * Helpers
 */

func checkStatusCode(expected int, rr *httptest.ResponseRecorder) error {
	if status := rr.Code; status != expected {
		return fmt.Errorf("Expected HTTP Status Code %d, got %d.  Details: %s", expected, status, rr.Body)
	}
	return nil
}

func checkHeader(header string, expected string, rr *httptest.ResponseRecorder) error {
	val := rr.Header().Get(header)
	if val != expected {
		return fmt.Errorf("Expected %s, got %s for HTTP header: %s", expected, val, header)
	}
	return nil
}

func checkResponseBody(expected interface{}, in interface{}, rr *httptest.ResponseRecorder) error {
	if err := json.NewDecoder(rr.Body).Decode(&in); err != nil {
		return fmt.Errorf("Invalid JSON response, details: %s", err)
	}
	if diff := deep.Equal(in, expected); diff != nil {
		return fmt.Errorf("%s", diff)
	}
	return nil
}
