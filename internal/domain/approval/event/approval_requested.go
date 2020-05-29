package event

import (
	"encoding/json"

	"forge.lmig.com/n1505471/pizza-shop/eventsource"
)

func init() {
	eventsource.RegisterEventType(&ApprovalRequested{})
}

// ApprovalRequested fired when approval is issued to the vendor system
type ApprovalRequested struct {
	ApprovalID int `json:"approvalId"`
}

func (e *ApprovalRequested) Version() int {
	return 1
}

func (e *ApprovalRequested) Load(data json.RawMessage, version int) error {
	switch version {
	default:
		err := json.Unmarshal(data, e)
		if err != nil {
			return err
		}
	}
	return nil
}
