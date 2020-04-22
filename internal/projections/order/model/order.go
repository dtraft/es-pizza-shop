package model

import "forge.lmig.com/n1505471/pizza-shop/internal/domain/order/model"

type Order struct {
	OrderID     string
	ServiceType model.ServiceType
	Description string
}
