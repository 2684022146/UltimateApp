package controller

import (
	"net/http"
	"strconv"
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
func (controller *OrderController) OrderDetailBasic(ctx *gin.Context) {
	orderNo := ctx.Query("order_no")
	c := ctx.Request.Context()
	order, err := controller.service.OrderDetailBasic(c, orderNo)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取订单基本详情失败")
		return
	}
	util.Success(ctx, order)
}

func (controller *OrderController) SenderFinishedOrder(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	var pageInfo model.Page
	if err := ctx.ShouldBindQuery(&pageInfo); err != nil {
		pageInfo.CurrentPage = 1
		pageInfo.PerPage = 10
	}
	// 确保分页参数有效
	if pageInfo.CurrentPage < 1 {
		pageInfo.CurrentPage = 1
	}
	if pageInfo.PerPage < 1 {
		pageInfo.PerPage = 10
	}
	c := ctx.Request.Context()
	orderList, total, err := controller.service.SenderFinishedOrder(c, userId, &pageInfo)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取寄件人已完成订单失败")
		return
	}
	util.SelectOrderSuccess(ctx, orderList, total)
}

func (controller *OrderController) SenderInTransitOrder(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	var pageInfo model.Page
	if err := ctx.ShouldBindQuery(&pageInfo); err != nil {
		pageInfo.CurrentPage = 1
		pageInfo.PerPage = 10
	}
	// 确保分页参数有效
	if pageInfo.CurrentPage < 1 {
		pageInfo.CurrentPage = 1
	}
	if pageInfo.PerPage < 1 {
		pageInfo.PerPage = 10
	}
	c := ctx.Request.Context()
	orderList, total, err := controller.service.SenderInTransitOrder(c, userId, &pageInfo)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取寄件人在途中订单失败")
		return
	}
	util.SelectOrderSuccess(ctx, orderList, total)
}
func (controller *OrderController) SenderWaitingOrder(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	var pageInfo model.Page
	if err := ctx.ShouldBindQuery(&pageInfo); err != nil {
		pageInfo.CurrentPage = 1
		pageInfo.PerPage = 10
	}
	// 确保分页参数有效
	if pageInfo.CurrentPage < 1 {
		pageInfo.CurrentPage = 1
	}
	if pageInfo.PerPage < 1 {
		pageInfo.PerPage = 10
	}
	c := ctx.Request.Context()
	orderList, total, err := controller.service.SenderWaitingOrder(c, userId, &pageInfo)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取待接单订单失败")
		return
	}
	util.SelectOrderSuccess(ctx, orderList, total)
}
func (controller *OrderController) SenderCencelOrder(ctx *gin.Context) {
	orderNo := ctx.Query("order_no")
	c := ctx.Request.Context()
	if err := controller.service.SenderCencelOrder(c, orderNo); err != nil {
		util.Fail(ctx, http.StatusInternalServerError, "取消订单失败")
		return
	}
	util.Success(ctx, nil)
}

// /////////////////////
func (controller *OrderController) ReceiverFinishedOrder(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	var pageInfo model.Page
	if err := ctx.ShouldBindQuery(&pageInfo); err != nil {
		pageInfo.CurrentPage = 1
		pageInfo.PerPage = 10
	}
	// 确保分页参数有效
	if pageInfo.CurrentPage < 1 {
		pageInfo.CurrentPage = 1
	}
	if pageInfo.PerPage < 1 {
		pageInfo.PerPage = 10
	}
	c := ctx.Request.Context()
	orderList, total, err := controller.service.ReceiverFinishedOrder(c, userId, &pageInfo)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取收件人已完成订单失败")
		return
	}
	util.SelectOrderSuccess(ctx, orderList, total)
}
func (controller *OrderController) ReceiverInTransitOrder(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	var pageInfo model.Page
	if err := ctx.ShouldBindQuery(&pageInfo); err != nil {
		pageInfo.CurrentPage = 1
		pageInfo.PerPage = 10
	}
	// 确保分页参数有效
	if pageInfo.CurrentPage < 1 {
		pageInfo.CurrentPage = 1
	}
	if pageInfo.PerPage < 1 {
		pageInfo.PerPage = 10
	}
	c := ctx.Request.Context()
	orderList, total, err := controller.service.ReceiverInTransitOrder(c, userId, &pageInfo)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取收件人途中订单失败")
		return
	}
	util.SelectOrderSuccess(ctx, orderList, total)
}

// ///////////////////
func (controller *OrderController) RiderOrderList(ctx *gin.Context) {
	var pageInfo model.Page
	statusStr := ctx.Query("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "状态参数错误")
		return
	}
	if err = ctx.ShouldBindQuery(&pageInfo); err != nil {
		pageInfo.CurrentPage = 1
		pageInfo.PerPage = 10
	}
	// 确保分页参数有效
	if pageInfo.CurrentPage < 1 {
		pageInfo.CurrentPage = 1
	}
	if pageInfo.PerPage < 1 {
		pageInfo.PerPage = 10
	}
	c := ctx.Request.Context()
	orderList, total, err := controller.service.RiderOrderList(c, status, &pageInfo)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取待接单订单失败")
		return
	}
	util.SelectOrderSuccess(ctx, orderList, total)
}

// 骑手接单
func (controller *OrderController) AcceptOrder(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	orderNo := ctx.Query("order_no")
	c := ctx.Request.Context()
	if err := controller.service.AcceptOrder(c, orderNo, userId); err != nil {
		util.Fail(ctx, http.StatusInternalServerError, "接单失败")
		return
	}
	util.Success(ctx, nil)
}

// 开始配送
func (controller *OrderController) StartDelivery(ctx *gin.Context) {
	orderNo := ctx.Query("order_no")
	c := ctx.Request.Context()
	if err := controller.service.StartDelivery(c, orderNo); err != nil {
		util.Fail(ctx, http.StatusInternalServerError, "开始配送失败")
		return
	}
	util.Success(ctx, nil)
}

// 取件
func (controller *OrderController) PickupOrder(ctx *gin.Context) {
	orderNo := ctx.Query("order_no")
	c := ctx.Request.Context()
	if err := controller.service.PickupOrder(c, orderNo); err != nil {
		util.Fail(ctx, http.StatusInternalServerError, "取件失败")
		return
	}
	util.Success(ctx, nil)
}

// 上传位置
func (controller *OrderController) UploadLocation(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	var req struct {
		OrderID   uint    `json:"order_id" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
		Latitude  float64 `json:"latitude" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	c := ctx.Request.Context()
	if err := controller.service.UploadLocation(c, req.OrderID, userId, req.Longitude, req.Latitude); err != nil {
		util.Fail(ctx, http.StatusInternalServerError, "上传位置失败")
		return
	}
	util.Success(ctx, nil)
}

// 完成订单
func (controller *OrderController) CompleteOrder(ctx *gin.Context) {
	orderNo := ctx.Query("order_no")
	c := ctx.Request.Context()
	if err := controller.service.CompleteOrder(c, orderNo); err != nil {
		util.Fail(ctx, http.StatusInternalServerError, "完成订单失败")
		return
	}
	util.Success(ctx, nil)
}
