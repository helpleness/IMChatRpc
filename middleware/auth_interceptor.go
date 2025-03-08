// 文件: middleware/grpc_interceptor.go
package middleware

import (
	"context"
	"github.com/helpleness/IMChatRpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

func UnaryAuthInterceptor(authService *auth.GRPCAuth) grpc.UnaryServerInterceptor {
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

		// 详细的错误处理和日志
		userID, err := authService.ValidateToken(ctx)
		if err != nil {
			log.Printf("认证失败 - 方法: %s, 错误: %v",
				info.FullMethod, err)

			// 根据错误类型返回不同的状态码
			var statusCode codes.Code
			switch err {
			case auth.ErrMissingToken:
				statusCode = codes.Unauthenticated
			case auth.ErrInvalidToken:
				statusCode = codes.PermissionDenied
			default:
				statusCode = codes.Internal
			}

			return nil, status.Errorf(statusCode, "认证失败: %v", err)
		}

		// 将用户ID添加到上下文
		newCtx := context.WithValue(ctx, "userid", userID)
		return handler(newCtx, req)
	}
}

// StreamAuthInterceptor 是一个流式RPC拦截器，用于验证流请求中的token
func StreamAuthInterceptor(authService *auth.GRPCAuth) grpc.StreamServerInterceptor {
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

		// 从流的上下文中获取元数据
		ctx := ss.Context()

		// 详细的元数据日志
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Println("无法获取流的元数据")
			return status.Errorf(codes.Unauthenticated, "未找到认证元数据")
		}

		// 打印所有元数据
		log.Printf("流式RPC元数据: %+v", md)

		// 尝试验证 Token
		userID, err := authService.ValidateToken(ctx)
		if err != nil {
			// 详细的错误日志
			log.Printf("流式RPC认证失败 - 方法: %s, 错误: %v",
				info.FullMethod, err)

			// 根据错误类型返回不同的状态码
			var statusCode codes.Code
			switch err {
			case auth.ErrMissingToken:
				statusCode = codes.Unauthenticated
			case auth.ErrInvalidToken:
				statusCode = codes.PermissionDenied
			default:
				statusCode = codes.Internal
			}

			return status.Errorf(statusCode, "流式RPC认证失败: %v", err)
		}

		// 将用户ID添加到上下文
		newCtx := context.WithValue(ctx, "userid", userID)

		// 包装流以传递修改后的上下文
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		// 日志记录成功认证
		log.Printf("流式RPC认证成功 - 用户ID: %d, 方法: %s",
			userID, info.FullMethod)

		return handler(srv, wrappedStream)
	}
}

// wrappedServerStream 包装了grpc.ServerStream以提供修改后的上下文
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 返回包装的上下文
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
