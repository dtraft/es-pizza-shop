package eventsource

import "fmt"

type AggregateLockError struct {
	ID       string
	Sequence int
}

func (err *AggregateLockError) Error() string {
	return fmt.Sprintf("AggregateLockError: Aggregate with id %s has already processed sequence: %d", err.ID, err.Sequence)
}
