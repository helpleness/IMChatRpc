// middleware/auth_interceptor.go
package middleware

import (
	"context"
	"github.com/helpleness/IMChatRpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthInterceptor 是一个 gRPC 拦截器，用于验证请求中的 token
func AuthInterceptor(validator auth.TokenValidator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// 健康检查不需要认证
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		token, err := auth.AuthFromMD(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		userID, err := validator.ValidateToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// 将用户ID添加到上下文
		newCtx := context.WithValue(ctx, "user_id", userID)
		return handler(newCtx, req)
	}
}

// StreamAuthInterceptor 是一个流式 gRPC 拦截器，用于验证流请求中的 token
func StreamAuthInterceptor(validator auth.TokenValidator) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// 健康检查不需要认证
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(srv, ss)
		}

		ctx := ss.Context()
		token, err := auth.AuthFromMD(ctx)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		userID, err := validator.ValidateToken(token)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// 包装ServerStream以传递修改后的上下文
		newCtx := context.WithValue(ctx, "user_id", userID)
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}
		return handler(srv, wrappedStream)
	}
}

// wrappedServerStream 包装了 grpc.ServerStream 以使用修改后的上下文
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 返回包装的上下文
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
