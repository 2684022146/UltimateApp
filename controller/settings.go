package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"webdemo/consts"
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
	userId_ctx, _ := ctx.Get("user_id")
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
	util.Success(ctx, nil)
}
func (controller *SettingsController) AddressList(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
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
func (controller *SettingsController) AddressDetail(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	addressIdStr := ctx.Query("address_id")
	if addressIdStr == "" {
		util.Fail(ctx, http.StatusBadRequest, consts.ParmFalse)
		return
	}
	addressId, err := strconv.ParseUint(addressIdStr, 10, 64)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, consts.TransformFalse)
		return
	}
	c := ctx.Request.Context()
	addressInfo, err := controller.settingsService.AddressDetail(c, uint(addressId), userId)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "获取地址详情失败")
		log.Println(err)
		return
	}
	util.Success(ctx, addressInfo)
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
func (controller *SettingsController) DeleteAddress(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	addressIdStr := ctx.Query("address_id")
	addressId, err := strconv.ParseUint(addressIdStr, 10, 64)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, consts.TransformFalse)
		log.Println(err, "addressIdStr", addressIdStr)
		return
	}
	c := ctx.Request.Context()
	err = controller.settingsService.DeleteAddress(c, uint(addressId), userId)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "删除地址失败")
		return
	}
	util.Success(ctx, nil)
}
func (controller *SettingsController) SetDefault(ctx *gin.Context) {
	userId_ctx, _ := ctx.Get("user_id")
	userId := userId_ctx.(uint)
	addressIdStr := ctx.Query("address_id")
	addressId, err := strconv.ParseUint(addressIdStr, 10, 64)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, consts.TransformFalse)
		log.Println(err, "addressIdStr", addressIdStr)
		return
	}
	c := ctx.Request.Context()
	err = controller.settingsService.SetDefault(c, uint(addressId), userId)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, "设为默认地址失败")
		return
	}
	util.Success(ctx, nil)
}
