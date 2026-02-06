package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1. 定义统一响应结构体（和前端约定的格式）
// 所有接口都返回该格式，前端只需解析这三个字段即可
type ApiResponse struct {
	Code int         `json:"code"` // 业务状态码（200成功，非200失败）
	Msg  string      `json:"msg"`  // 响应信息（成功/失败提示）
	Data interface{} `json:"data"` // 响应数据（成功时返回业务数据，失败时为nil）
}

// 2. Success 成功响应工具方法
// 入参：gin.Context（用于返回响应）、data（需要返回给前端的业务数据）
func Success(c *gin.Context, data interface{}) {
	// 构造统一响应结构体
	response := ApiResponse{
		Code: 200,    // 成功固定业务状态码200
		Msg:  "操作成功", // 默认成功提示（可根据需求修改）
		Data: data,   // 传入的业务数据（如登录成功的token）
	}

	// 返回JSON响应，HTTP状态码固定为200（业务状态码在response.Code中区分）
	c.JSON(http.StatusOK, response)
}

// 3. Fail 失败响应工具方法
// 入参：gin.Context、业务状态码、失败提示信息
func Fail(c *gin.Context, code int, msg string) {
	// 构造统一响应结构体
	response := ApiResponse{
		Code: code, // 传入的业务失败码（如400参数错误、401未授权）
		Msg:  msg,  // 传入的失败提示信息（如"用户名或密码错误"）
		Data: nil,  // 失败时数据为nil
	}

	// 返回JSON响应，HTTP状态码统一为200（也可根据需求改为对应状态码，如400、500）
	// 推荐统一返回200，避免前端处理跨域和状态码判断的额外开销，通过业务码区分错误
	c.JSON(http.StatusOK, response)
}

// 4. 可选：扩展方法——带自定义成功提示的Success
// 适用于需要自定义成功提示的场景（如"登录成功"、"查询成功"）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	response := ApiResponse{
		Code: 200,
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}
