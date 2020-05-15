package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	domain "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
	"forge.lmig.com/n1505471/pizza-shop/internal/projections/order/model"
	"forge.lmig.com/n1505471/pizza-shop/internal/projections/order/repository"
	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/julienschmidt/httprouter"
)

var svc = dynamodb.New(session.New(), aws.NewConfig())
var repo = repository.NewRepository(svc, os.Getenv("TABLE_NAME"))
var router = httprouter.New()

func init() {
	router.GET("/orders", queryAllOrders)
	router.GET("/orders/:orderID", getOrder)
}

func main() {
	log.Fatal(gateway.ListenAndServe(":3000", router))
}

/*
 * Routes
 */

func queryAllOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	orders, err := repo.QueryAllOrders()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var resources []*orderResource
	for _, o := range orders {
		resources = append(resources, resourceFromOrder(o))
	}

	log.Printf("Orders: %+v", orders)

	jsonResponse(w, resources)
}

func getOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	orderID := p.ByName("orderID")

	order, err := repo.GetOrder(orderID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Printf("Order: %+v", order)

	jsonResponse(w, resourceFromOrder(order))
}

/*
 * Resources
 */

type orderResource struct {
	OrderID     string             `json:"orderId"`
	ServiceType domain.ServiceType `json:"serviceType"`
	Status      domain.Status      `json:"status"`
	Description string             `json:"description"`

	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

func resourceFromOrder(o *model.Order) *orderResource {
	status := domain.Started
	if o.Status > 0 {
		status = o.Status
	}

	return &orderResource{
		OrderID:     o.OrderID,
		ServiceType: o.ServiceType,
		Status:      status,
		Description: o.Description,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

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
