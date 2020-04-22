package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
	router.POST("/orders/:orderID/toggle", c.toggleOrder)
	router.POST("/orders", c.startOrder)
}

func init() {
	var svc *dynamodb.DynamoDB
	f := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if strings.Contains(f, "local") {
		svc = dynamodb.New(session.New(), aws.NewConfig().WithRegion("localhost").WithEndpoint("http://host.docker.internal:9898"))
	} else {
		svc = dynamodb.New(session.New(), aws.NewConfig())
	}

	store := ddbES.New(svc, os.Getenv("TABLE_NAME"))
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

func (c *Controller) toggleOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	orderID := p.ByName("orderID")

	if err := c.orderSvc.ToggleOrderServiceType(orderID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Order toggled.")
	return
}

func (c *Controller) startOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var resource *orderResource
	err := json.NewDecoder(r.Body).Decode(&resource)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(resource); err != nil {
		invalidResponse(w, err)
		return
	}

	orderID, err := c.orderSvc.StartOrder(resource.toOrder())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resource.OrderID = orderID
	jsonResponse(w, resource)
}

/*
 * Types
 */

type orderResource struct {
	OrderID     string            `json:"orderId"`
	ServiceType model.ServiceType `json:"serviceType" validate:"gte=1,lte=2"`
	Description string            `json:"description"`
}

func (o *orderResource) toOrder() *model.Order {
	return &model.Order{
		OrderID:     o.OrderID,
		ServiceType: o.ServiceType,
		Description: o.Description,
	}
}

// func (o *OrderResource) validate

/*
 * Helpers
 */

func jsonResponse(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	resp := &errorResponse{
		OK:     false,
		Errors: errors,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

type errorResponse struct {
	OK     bool `json:"ok"`
	Errors []*validationError
}

type validationError struct {
	Field   string
	Message string
}
