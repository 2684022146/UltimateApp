package controller

import (
	"net/http"
	"webdemo/model"
	"webdemo/service"
	"webdemo/util"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginService service.LoginService
}

func NewLoginController(loginService service.LoginService) *LoginController {
	return &LoginController{
		loginService: loginService,
	}
}
func (controller *LoginController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "param form error",
			"data": nil,
		})
		return
	}
	c := ctx.Request.Context()
	token, err := controller.loginService.Login(c, &req)
	if err != nil {
		util.Fail(ctx, http.StatusBadGateway, err.Error())
		return
	}
	util.Success(ctx, gin.H{
		"token": token,
	})
}
func (controller *LoginController) Regist(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	c := ctx.Request.Context()
	token, err := controller.loginService.Regist(c, &req)
	if err != nil {
		util.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	util.Success(ctx, gin.H{
		"token": token,
	})
}
