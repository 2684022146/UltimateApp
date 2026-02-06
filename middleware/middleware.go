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
		ctx.Set("userId", claims.ID)
		ctx.Set("username", claims.Username)
		ctx.Set("role_id", claims.RoleID)
		ctx.Next()
	}

}
