package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"webdemo/service"
	"webdemo/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	//全局服务权限实例
	permissionService service.PermissionService
)

// 设置权限服务实例
func SetPermissionService(ps service.PermissionService) {
	permissionService = ps
}
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			util.Fail(ctx, http.StatusUnauthorized, "请先登陆")
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.Fail(ctx, http.StatusUnauthorized, "Token格式错误")
			ctx.Abort()
			return
		}
		tokenString := parts[1]
		claims := &util.CustomClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			//验证签名算法是否一致
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			//返回jwt密钥
			return []byte(util.JwtSecret), nil
		})
		//处理验证错误
		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenExpired) {
				util.Fail(ctx, http.StatusUnauthorized, "Token已过期")
				ctx.Abort()
				return
			}
			util.Fail(ctx, http.StatusUnauthorized, "Token错误")
			ctx.Abort()
			return
		}
		ctx.Set("user_id", claims.UserId)
		ctx.Set("username", claims.Username)
		ctx.Set("role_id", claims.RoleID)
		ctx.Set("token", tokenString)
		ctx.Next()
	}

}

// 角色权限校验中间件
func RequireRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取当前角色
		roleId, exists := ctx.Get("role_id")
		if !exists {
			util.Fail(ctx, http.StatusUnauthorized, "获取角色失败")
			ctx.Abort()
			return
		}
		//获取当前请求api
		method := ctx.Request.Method
		apiPath := ctx.Request.URL.Path
		//角色权限校验
		hasPermission := permissionService.CheckPermission(roleId.(int8), apiPath, method)
		if !hasPermission {
			log.Println("hasPermission", hasPermission)
			util.Fail(ctx, http.StatusUnauthorized, "权限不足")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
