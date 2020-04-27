package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/markphelps/optional"

	"github.com/go-playground/validator/v10"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	ddbES "forge.lmig.com/n1505471/pizza-shop/eventsource/store/dynamodb"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"
	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/julienschmidt/httprouter"
)

var router *httprouter.Router
var validate *validator.Validate

var controller *Controller

type Controller struct {
	orderSvc order.ServiceAPI
}

func (c *Controller) registerRoutes(router *httprouter.Router) {
	router.POST("/orders", c.startOrder)
	router.PATCH("/orders/:orderID", c.updateOrder)
}

func init() {
	var store eventsource.EventStore
	f := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if strings.Contains(f, "local") {
		svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion("localhost").WithEndpoint("http://host.docker.internal:9898"))
		store = ddbES.New(svc, "EventsTable-local")
	} else {
		svc := dynamodb.New(session.New(), aws.NewConfig())
		store = ddbES.New(svc, os.Getenv("TABLE_NAME"))
	}

	es := eventsource.New(store)
	orderSvc := order.NewService(es)

	controller = &Controller{
		orderSvc: orderSvc,
	}

	validate = validator.New()

	router = httprouter.New()
	controller.registerRoutes(router)
}

func main() {
	log.Fatal(gateway.ListenAndServe(":3000", router))
}

func (c *Controller) startOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var resource *orderResource
	err := json.NewDecoder(r.Body).Decode(&resource)

	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	if err := validate.Struct(resource); err != nil {
		invalidResponse(w, err)
		return
	}

	orderID, err := c.orderSvc.StartOrder(resource.toOrder())
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	resource.OrderID = orderID
	jsonResponse(w, &response{
		OK: true,
		Result: map[string]string{
			"orderId": orderID,
		},
	})
}

func (c *Controller) updateOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	orderID := p.ByName("orderID")

	var resource *orderPatchResource
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	resource.OrderID = orderID

	if err := c.orderSvc.UpdateOrder(resource.toOrderPatch()); err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	jsonResponse(w, &response{
		OK: true,
	})
}

/*
 * Types
 */

type orderResource struct {
	OrderID     string            `json:"orderId"`
	ServiceType model.ServiceType `json:"serviceType"`
	Description string            `json:"description"`
}

func (o *orderResource) toOrder() *model.Order {
	return &model.Order{
		OrderID:     o.OrderID,
		ServiceType: o.ServiceType,
		Description: o.Description,
	}
}

type orderPatchResource struct {
	OrderID     string                    `json:"orderId"`
	ServiceType model.OptionalServiceType `json:"serviceType"`
	Description optional.String           `json:"description"`
}

func (o *orderPatchResource) toOrderPatch() *model.OrderPatch {
	return &model.OrderPatch{
		OrderID:     o.OrderID,
		ServiceType: o.ServiceType,
		Description: o.Description,
	}
}

/*
 * Helpers
 */

func jsonResponse(w http.ResponseWriter, body interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func invalidResponse(w http.ResponseWriter, err error) {
	var errors []*validationError
	for _, err := range err.(validator.ValidationErrors) {
		var message = fmt.Sprintf("%s is not valid for field %s", err.Value(), err.Field())
		if err.Tag() == "required" {
			message = fmt.Sprintf("%s is a required field.", err.Field())
		}
		errors = append(errors, &validationError{
			Field:   err.Namespace(),
			Message: message,
		})
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	resp := &response{
		OK:     false,
		Errors: errors,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func errorResponse(w http.ResponseWriter, err error, code int) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)

	resp := &response{
		OK:     false,
		Result: err.Error(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

type response struct {
	OK     bool               `json:"ok"`
	Result interface{}        `json:"result,omitempty"`
	Errors []*validationError `json:"errors,omitempty"`
}

type validationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
