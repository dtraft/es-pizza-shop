package command

import "strconv"

// ReceiveApproval attempts to submit the order for fulfillment
type RequestDelivery struct {
	DeliveryID int
}

func (c *RequestDelivery) AggregateID() string {
	return strconv.Itoa(c.DeliveryID)
}
