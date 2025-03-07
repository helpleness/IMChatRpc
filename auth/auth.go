// auth/auth.go
package auth

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

var (
	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New("invalid token")
)

// TokenValidator 定义了验证token的接口
type TokenValidator interface {
	ValidateToken(token string) (string, error)
}

// SimpleTokenValidator 是一个简单的token验证器
type SimpleTokenValidator struct {
	// 在实际应用中，你可能会有一个token到用户ID的映射，或使用JWT
	validTokens map[string]string
}

// NewSimpleTokenValidator 创建一个新的SimpleTokenValidator
func NewSimpleTokenValidator() *SimpleTokenValidator {
	// 这只是示例，实际中你会使用更安全的方式存储和验证token
	return &SimpleTokenValidator{
		validTokens: map[string]string{
			"token123": "user1",
			"token456": "user2",
		},
	}
}

// ValidateToken 检查token是否有效，并返回关联的用户ID
func (v *SimpleTokenValidator) ValidateToken(token string) (string, error) {
	userID, exists := v.validTokens[token]
	if !exists {
		return "", ErrInvalidToken
	}
	return userID, nil
}

// GetUserIDFromContext 从上下文中提取用户ID
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}

// AuthFromMD 从元数据中提取token
func AuthFromMD(ctx context.Context) (string, error) {
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
