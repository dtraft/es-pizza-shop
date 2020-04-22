package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

	"github.com/julienschmidt/httprouter"
)

type Condition func(rr *httptest.ResponseRecorder) error

func TestStartOrder(t *testing.T) {

	malformedReq, _ := http.NewRequest("POST", "/orders", strings.NewReader(`
			"serviceType": 0,
			"description": "Here is a description."
		}
	`))

	invalidReq, _ := http.NewRequest("POST", "/orders", strings.NewReader(`
		{
			"serviceType": 0,
			"description": "Here is a description."
		}
	`))

	validReq, _ := http.NewRequest("POST", "/orders", strings.NewReader(`
		{
			"serviceType": 1,
			"description": "Here is a description."
		}
	`))

	validReq2, _ := http.NewRequest("POST", "/orders", strings.NewReader(`
		{
			"serviceType": 1,
			"description": "Here is a description."
		}
	`))

	cases := []struct {
		svc       order.ServiceAPI
		request   *http.Request
		condition Condition
	}{
		{
			svc:     &mockOrderService{},
			request: malformedReq,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusBadRequest, rr); err != nil {
					return err
				}
				return nil
			},
		},
		{
			svc:     &mockOrderService{},
			request: invalidReq,
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
			request: validReq,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusBadRequest, rr); err != nil {
					return err
				}
				return nil
			},
		},
		{
			svc:     &mockOrderService{},
			request: validReq2,
			condition: func(rr *httptest.ResponseRecorder) error {
				if err := checkStatusCode(http.StatusOK, rr); err != nil {
					return err
				}

				if err := checkHeader("content-type", "application/json", rr); err != nil {
					return err
				}

				expected := &orderResource{
					OrderID:     "orderId",
					ServiceType: 1,
					Description: "Here is a description.",
				}
				if err := checkResponseBody(expected, rr); err != nil {
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

func checkResponseBody(expected interface{}, rr *httptest.ResponseRecorder) error {
	rawType := reflect.TypeOf(expected)
	// source is a pointer, convert to its value
	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}
	resource := reflect.New(rawType).Interface()
	if err := json.NewDecoder(rr.Body).Decode(&resource); err != nil {
		return fmt.Errorf("Invalid JSON response, details: %s", err)
	}
	if diff := deep.Equal(resource, expected); diff != nil {
		return fmt.Errorf("%s", diff)
	}
	return nil
}
