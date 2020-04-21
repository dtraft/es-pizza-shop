default: clean infrastructure_eventforwarder order_writeapi order_readapi order_projection

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

# Tests
test:
	go test ./... -cover -race -coverprofile=coverage.out

# Utils
clean: 
	rm -rf .bin/*