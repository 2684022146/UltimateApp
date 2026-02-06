package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JwtSecret = "abcd-efgh-ijkl-mnop"

type CustomClaim struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleID   int8   `json:"role_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId uint, username string, roleId int8, expireDuration time.Duration) (string, error) {
	if userId <= 0 || username == "" || expireDuration <= 0 {
		return "", fmt.Errorf("generate token fail")
	}
	now := time.Now()
	customeClaims := &CustomClaim{
		UserId:   userId,
		Username: username,
		RoleID:   roleId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "yzx",
			Subject:   fmt.Sprintf("%d", userId),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customeClaims)
	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", fmt.Errorf("signature fail:%w", err)
	}
	return tokenString, nil
}
func Md5String(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
