package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

func init() {
	eventsource.RegisterEventType(&ApprovalReceived{})
}

// ApprovalReceived fired when approval is received from vendor system
type ApprovalReceived struct {
	ApprovalID int `json:"approvalId"`
}

func (e *ApprovalReceived) Version() int {
	return 1
}

func (e *ApprovalReceived) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}
