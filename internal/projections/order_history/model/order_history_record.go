package model

import (
	"time"
)

type OrderHistoryRecord struct {
	OrderID         string     `json:"orderId,omitempty"`
	Description     string     `json:"transactionDescription,omitempty"`
	TransactionDate *time.Time `json:"transactionDate,omitempty"`
}
