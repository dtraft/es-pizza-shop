default: clean infrastructure_eventforwarder order_writeapi order_readapi order_projection order_fulfillment_saga

# Local Dev
local:
	npm run start

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

# Replays
order_projection_replay:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/order_projection_replay lambda/order/replay/order_projection_replay.go

order_history_projection_replay:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/order_history_projection_replay lambda/order/replay/order_history_projection_replay.go

package_order_projection_replay: order_projection_replay
	docker build .bin -f lambda/order/replay.dockerfile -t order-projection-replay

package_order_history_projection_replay: order_history_projection_replay
	docker build .bin -f lambda/order/history_replay.dockerfile -t order-history-projection-replay

# Tests
test:
	go test ./... -cover -race -coverprofile=coverage.out

# Coverage
coverage: test
	go tool cover -html=coverage.out

# Utils
clean: 
	rm -rf .bin/*