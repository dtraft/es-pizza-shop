package model

// ServiceType is the type of order, e.g. pickup or delivery
type ServiceType int

const (
	Pickup ServiceType = iota
	Delivery
)
