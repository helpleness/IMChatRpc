package controller

import (
	"context"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
)

// HealthService 实现 grpc_health_v1.HealthServer 接口
type HealthService struct {
	mu     sync.Mutex
	status map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
}

// NewHealthService 创建新的 HealthService 实例
func NewHealthService() *HealthService {
	return &HealthService{
		status: make(map[string]grpc_health_v1.HealthCheckResponse_ServingStatus),
	}
}

// Check 方法实现
func (h *HealthService) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	service := req.Service
	status, ok := h.status[service]
	if !ok {
		status = grpc_health_v1.HealthCheckResponse_UNKNOWN
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: status,
	}, nil
}

// Watch 方法实现（非必须，可返回未实现错误）
func (h *HealthService) Watch(req *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return nil
}

// SetStatus 设置服务健康状态
func (h *HealthService) SetStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status[service] = status
}

// RegisterHealthService 将健康检查服务注册到 gRPC 服务
func RegisterHealthService(server *grpc.Server, healthService *HealthService) {
	grpc_health_v1.RegisterHealthServer(server, healthService)
	// 获取服务名
	Name := viper.GetString("consul.name")
	// 设置默认服务健康状态
	healthService.SetStatus("", grpc_health_v1.HealthCheckResponse_SERVING) // 默认服务
	healthService.SetStatus(Name, grpc_health_v1.HealthCheckResponse_SERVING)
}
