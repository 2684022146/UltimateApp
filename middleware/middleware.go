package middleware

import (
	"errors"
	"net/http"
	"strings"
	"webdemo/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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
		ctx.Set("user_id", claims.ID)
		ctx.Set("username", claims.Username)
		ctx.Set("role_id", claims.RoleID)
		ctx.Set("token", tokenString)
		ctx.Next()
	}

}
func RequireRole(roleId int8) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authMiddleware := AuthMiddleware()
		authMiddleware(ctx)
		if ctx.IsAborted() {
			return
		}
		//验证角色
		userRole, exists := ctx.Get("role_id")
		if !exists {
			util.Fail(ctx, http.StatusUnauthorized, "获取角色失败")
			ctx.Abort()
			return
		}
		if userRole.(int8) != roleId {
			util.Fail(ctx, http.StatusForbidden, "权限不足")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// RequireConsignee 收货人角色验证中间件
func RequireConsignee() gin.HandlerFunc {
	return RequireRole(0)
}

// RequireDeliveryman 配送员角色验证中间件
func RequireDeliveryman() gin.HandlerFunc {
	return RequireRole(1)
}
