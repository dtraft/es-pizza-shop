package model

// ServiceType is the type of order, e.g. pickup or delivery
type ServiceType int

const (
	_ ServiceType = iota
	Pickup
	Delivery
)

func (r ServiceType) String() string {
	return _ServiceTypeValueToName[r]
}

//go:generate jsonenums -type=ServiceType
//go:generate optional -type=ServiceType
