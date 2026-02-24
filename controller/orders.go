package controller

import (
	"fmt"
	"net/http"
	"webdemo/model"
	"webdemo/service"
	"webdemo/util"

	"github.com/gin-gonic/gin"
)

type OrdersController struct {
	ordersService service.OrdersService
}

func NewOrdersController(ordersService service.OrdersService) *OrdersController {
	return &OrdersController{
		ordersService: ordersService,
	}
}

func (controller *OrdersController) OrderList(ctx *gin.Context) {
	userId_ctx, exists := ctx.Get("user_id")
	if !exists {
		util.Fail(ctx, http.StatusUnauthorized, "请重新登陆")
		return
	}
	userId := userId_ctx.(uint)
	c := ctx.Request.Context()
	orders, err := controller.ordersService.OrderList(c, userId)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取订单列表失败")
		return
	}
	util.Success(ctx, orders)
}

func (controller *OrdersController) CreateOrder(ctx *gin.Context) {
	userId_ctx, exist := ctx.Get("user_id")
	if !exist {
		util.Fail(ctx, http.StatusUnauthorized, "请重新登录")
		return
	}
	var req *model.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	userId := userId_ctx.(uint)
	c := ctx.Request.Context()
	if err := controller.ordersService.CreateOrder(c, req, userId); err != nil {
		util.Fail(ctx, http.StatusBadRequest, fmt.Sprintf("%s", err))
		return
	}
	util.Success(ctx, nil)
}
