package model

// ServiceType is the type of order, e.g. pickup or delivery
type Status int

const (
	_ Status = iota
	Started
	Submitted
	Approved
	Delivered
)

func (r Status) String() string {
	return _StatusValueToName[r]
}

//go:generate jsonenums -type=Status
