package model

import (
	"time"

	"forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"
)

type Order struct {
	OrderID     string            `json:"orderId,omitempty"`
	ServiceType model.ServiceType `json:"serviceType,omitempty"`
	Description string            `json:"description,omitempty"`
	Status      model.Status      `json:"status,omitempty"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
