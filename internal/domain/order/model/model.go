package model

import "github.com/markphelps/optional"

type Order struct {
	OrderID     string
	ServiceType ServiceType
	Description string
}

type OrderPatch struct {
	OrderID     string
	ServiceType OptionalServiceType
	Description optional.String
}
