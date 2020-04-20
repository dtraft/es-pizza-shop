package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

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

	log.Printf("Orders: %+v", orders)

	jsonResponse(w, orders)
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