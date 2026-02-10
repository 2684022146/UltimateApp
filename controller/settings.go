package controller

import (
	"fmt"
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
	var userId int
	if !exist {
		util.Fail(ctx, http.StatusUnauthorized, "请重新登录")
	}
	var req *model.CreateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.Fail(ctx, http.StatusBadRequest, "参数错误")
	}
	c := ctx.Request.Context()
	v, ok := userId_ctx.(int)
	if ok {
		userId = v
	}
	if err := controller.settingsService.CreateAddress(c, req, userId); err != nil {
		util.Fail(ctx, http.StatusBadRequest, fmt.Sprintf("%s", err))
		return
	}
	util.Success(ctx, "success")
}
