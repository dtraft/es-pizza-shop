package main

import (
	"log"
	"os"

	"forge.lmig.com/n1505471/pizza-shop/eventsource/store/s3Store"

	"forge.lmig.com/n1505471/pizza-shop/internal/projections/order_history/repository"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	history "forge.lmig.com/n1505471/pizza-shop/internal/projections/order_history"

	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	log.Printf("Starting replay!")

	s3Svc := s3.New(session.New())
	ddbSvc := dynamodb.New(session.New())

	svc := s3Store.New(s3Svc, os.Getenv("EVENT_BUCKET"))

	repo := repository.NewRepository(ddbSvc, os.Getenv("TABLE_NAME"))
	projection := history.NewProjection(repo)

	pagecount := 0
	// List events in s3Store, page by page
	for {
		pagecount++
		log.Printf("-----Processing page %d -----", pagecount)
		events, err := svc.GetNextPage()
		if err != nil {
			log.Fatal(err)
		}

		for _, event := range events {
			// Apply projection
			log.Printf("Processing event: %+v", event)
			if err := projection.HandleEvent(event); err != nil {
				log.Fatal(err)
			}
		}

		if !svc.HasNextPage() {
			break
		}

	}

	log.Printf("Replay is complete!")
}
