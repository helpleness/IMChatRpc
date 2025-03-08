// 文件: auth/grpc_auth.go
package auth

import (
	"context"
	"errors"
	"github.com/helpleness/IMChatRpc/database"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
	"time"
	// 导入你现有的中间件包，获取JWT解析函数
	Middleware "github.com/helpleness/IMChatAdmin/middleware"
)

var (
	ErrMissingToken = errors.New("缺少认证令牌")
	ErrInvalidToken = errors.New("无效的认证令牌")
)

// GRPCAuth 提供gRPC认证服务
type GRPCAuth struct {
	// 可以添加额外依赖，例如Redis客户端或数据库连接
}

// NewGRPCAuth 创建新的gRPC认证服务
func NewGRPCAuth() *GRPCAuth {
	return &GRPCAuth{}
}

// AuthFromMetadata 从gRPC元数据中提取token
func AuthFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingToken
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", ErrMissingToken
	}

	// 解析 "Bearer token" 格式
	auth := values[0]
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", ErrInvalidToken
	}
	return strings.TrimPrefix(auth, "Bearer "), nil
}

func (a *GRPCAuth) ValidateToken(ctx context.Context) (uint, error) {
	// 日志记录元数据信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("无法获取元数据")
		return 0, ErrMissingToken
	}
	log.Printf("收到的元数据: %v", md)

	// 提取 Token
	tokenString, err := AuthFromMetadata(ctx)
	if err != nil {
		log.Printf("提取 Token 失败: %v", err)
		return 0, err
	}

	// 日志记录 Token 信息
	log.Printf("收到的 Token: %s", tokenString)

	// 解析 Token
	_, claims, err := Middleware.ParseToken(tokenString)
	if err != nil {
		log.Printf("Token 解析失败: %v", err)
		return 0, ErrInvalidToken
	}

	// 详细的 Token 验证日志
	log.Printf("Token 信息: UserID=%v, ExpiresAt=%v",
		claims.UserID,
		claims.ExpiresAt)

	// 过期检查
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		log.Printf("Token 已过期: %v", claims.ExpiresAt)
		return 0, ErrInvalidToken
	}

	// 用户存在性检查
	rediscli := database.GetRedisClient()
	db := database.GetDB()
	_, err = Middleware.Isuserexist(ctx, int(claims.UserID), db, rediscli)
	if err != nil {
		log.Printf("用户存在性检查失败: %v", err)
		return 0, err
	}

	return claims.UserID, nil
}

// GetUserIDFromContext 从context中获取用户ID
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value("userid").(uint)
	return id, ok
}
