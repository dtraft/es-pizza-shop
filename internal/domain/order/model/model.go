package model

// ServiceType is the type of order, e.g. pickup or delivery
type ServiceType int

const (
	Pickup ServiceType = iota + 1
	Delivery
)

type Order struct {
	OrderID     string      `json:"orderId"`
	ServiceType ServiceType `json:"serviceType"`
	Description string      `json:"description"`
}
