package controller

import (
	"net/http"
	"webdemo/model"
	"webdemo/service"
	"webdemo/util"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service service.OrderService
}

func NewOrdersController(service service.OrderService) *OrderController {
	return &OrderController{
		service: service,
	}
}
func (controller *OrderController) CreateOrder(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	var req *model.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	c := ctx.Request.Context()
	if err := controller.service.CreateOrder(c, req, userId); err != nil {
		util.Fail(ctx, http.StatusInternalServerError, "创建订单失败")
		return
	}
	util.Success(ctx, nil)
}
