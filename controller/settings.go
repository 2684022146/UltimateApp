package controller

import (
	"fmt"
	"log"
	"net/http"
	"webdemo/model"
	"webdemo/service"
	"webdemo/util"

	"github.com/gin-gonic/gin"
)

type SettingsController struct {
	settingsService service.SettingsService
}

func NewSettingsController(settingsService service.SettingsService) *SettingsController {
	return &SettingsController{
		settingsService: settingsService,
	}
}
func (controller *SettingsController) CreateAddress(ctx *gin.Context) {
	userId_ctx, exist := ctx.Get("user_id")
	if !exist {
		util.Fail(ctx, http.StatusUnauthorized, "请重新登录")
	}
	var req *model.CreateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	userId := userId_ctx.(uint)
	c := ctx.Request.Context()
	if err := controller.settingsService.CreateAddress(c, req, userId); err != nil {
		util.Fail(ctx, http.StatusBadRequest, fmt.Sprintf("%s", err))
		return
	}
	util.Success(ctx, "success")
}
func (controller *SettingsController) AddressList(ctx *gin.Context) {
	userId_ctx, exists := ctx.Get("user_id")
	if !exists {
		util.Fail(ctx, http.StatusUnauthorized, "请重新登陆")
		return
	}
	userId := userId_ctx.(uint)
	c := ctx.Request.Context()
	addressSlice, err := controller.settingsService.AddressList(c, userId)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取地址列表失败")
		return
	}
	log.Println(addressSlice)
	util.Success(ctx, addressSlice)

}
func (controller *SettingsController) UpdateAddress(ctx *gin.Context) {
	var req *model.Address

	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	log.Println(req.IsDefault)
	c := ctx.Request.Context()
	if err := controller.settingsService.UpdateAddress(c, req); err != nil {
		util.Fail(ctx, http.StatusBadRequest, "修改地址失败")
		return
	}
	util.Success(ctx, "success")
}
