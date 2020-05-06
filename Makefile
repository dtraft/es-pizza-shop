default: clean infrastructure_eventforwarder order_writeapi order_readapi order_projection order_fulfillment_saga

# Infrastructure
infrastructure_eventforwarder:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/infrastructure_eventforwarder lambda/infrastructure/eventforwarder/eventforwarder.go

# Order
order_writeapi:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/order_writeapi lambda/order/api/writeapi/order_writeapi.go

order_readapi:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/order_readapi lambda/order/api/readapi/order_readapi.go

order_projection:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/order_projection lambda/order/projection/order_projection.go

order_fulfillment_saga:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/order_fulfillment_saga lambda/order/saga/order_fulfillment_saga.go

# Tests
test:
	go test ./... -cover -race -coverprofile=coverage.out

# Coverage
coverage: test
	go tool cover -html=coverage.out

# Local Dev
local:
	npm run start

# Utils
clean: 
	rm -rf .bin/*