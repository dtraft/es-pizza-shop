.PHONY: default api eventforwarder clean orderprojection

default: clean api eventforwarder orderprojection

api:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/api lambda/api/main.go

eventforwarder:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/eventforwarder lambda/infrastructure/eventforwarder.go

orderprojection:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/orderprojection lambda/projections/order.go

clean: 
	rm -rf .bin/*