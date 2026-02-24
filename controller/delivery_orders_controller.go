package controller

import (
	"net/http"

	"webdemo/service"
	"webdemo/util"

	"github.com/gin-gonic/gin"
)

type DeliveryOrdersController struct {
	deliveryOrdersService service.DeliveryOrdersService
}

func NewDeliveryOrdersController(deliveryOrdersService service.DeliveryOrdersService) *DeliveryOrdersController {
	return &DeliveryOrdersController{
		deliveryOrdersService: deliveryOrdersService,
	}
}

func (controller *DeliveryOrdersController) DeliveryOrderList(ctx *gin.Context) {
	deliveryUserId_ctx, exists := ctx.Get("user_id")
	if !exists {
		util.Fail(ctx, http.StatusUnauthorized, "请重新登陆")
		return
	}
	deliveryUserId := deliveryUserId_ctx.(uint)
	c := ctx.Request.Context()
	orders, err := controller.deliveryOrdersService.DeliveryOrderList(c, deliveryUserId)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取待配送订单列表失败")
		return
	}
	util.Success(ctx, orders)
}
