package controller

import (
	"webdemo/service"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController(service service.OrderService) *OrderController {
	return &OrderController{
		service: service,
	}
}
