package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
	ddbES "forge.lmig.com/n1505471/pizza-shop/eventsource/store/dynamodb"
	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order"
	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/julienschmidt/httprouter"
)

var orderSvc *order.Service

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

	orderSvc = order.NewService(es)
}

func main() {
	router := httprouter.New()
	router.GET("/toggle/:orderID", toggleOrder)
	router.GET("/start", startOrder)
	log.Fatal(gateway.ListenAndServe(":3000", router))
}

func toggleOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	orderID := p.ByName("orderID")

	if err := orderSvc.ToggleOrderServiceType(orderID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, "Order toggled.")
	return
}

func startOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	orderID, err := orderSvc.StartOrder()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "Order Started: $s", orderID)
	return
}
