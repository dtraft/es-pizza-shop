package model

import "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

type Order struct {
	OrderID     string            `json:"orderId"`
	ServiceType model.ServiceType `json:"serviceType"`
}
