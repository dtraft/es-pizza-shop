default: clean writeapi readapi eventforwarder orderprojection

writeapi:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/writeapi lambda/api/writeapi.go

readapi:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/readapi lambda/api/readapi.go

eventforwarder:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/eventforwarder lambda/infrastructure/eventforwarder.go

orderprojection:
	env GOOS=linux go build -ldflags="-s -w"  -o .bin/orderprojection lambda/projections/order.go

clean: 
	rm -rf .bin/*