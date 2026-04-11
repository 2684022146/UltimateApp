package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JwtSecret = "abcd-efgh-ijkl-mnop"

// 雪花算法参数
const (
	workerIDBits  = 10
	sequenceBits  = 12
	maxWorkerID   = -1 ^ (-1 << workerIDBits)
	maxSequence   = -1 ^ (-1 << sequenceBits)
	timestampLeft = workerIDBits + sequenceBits
	workerIDLeft  = sequenceBits
	// 2024-01-01 00:00:00 UTC
	epoch = 1704067200000
)

// Snowflake 雪花算法生成器
type Snowflake struct {
	mu       sync.Mutex
	workerID int64
	lastTime int64
	sequence int64
}

// NewSnowflake 创建雪花算法生成器
func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	return &Snowflake{
		workerID: workerID,
		lastTime: -1,
		sequence: 0,
	}, nil
}

// Generate 生成雪花ID
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < s.lastTime {
		panic("clock is moving backwards")
	}

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	return ((now - epoch) << timestampLeft) | (s.workerID << workerIDLeft) | s.sequence
}

// 全局雪花算法生成器
var snowflake *Snowflake

func init() {
	var err error
	snowflake, err = NewSnowflake(1)
	if err != nil {
		panic(err)
	}
}

// GenerateOrderNo 生成订单编号
func GenerateOrderNo() string {
	id := snowflake.Generate()
	return fmt.Sprintf("ORD%d", id)
}

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
