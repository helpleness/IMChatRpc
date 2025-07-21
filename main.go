package main

import (
	"context"
	"github.com/helpleness/IMChatRpc/auth"
	"github.com/helpleness/IMChatRpc/config"
	"github.com/helpleness/IMChatRpc/controller"
	"github.com/helpleness/IMChatRpc/database"
	"github.com/helpleness/IMChatRpc/middleware"
	"github.com/helpleness/IMChatRpc/router"
	"github.com/helpleness/IMChatRpc/service"
	pb "github.com/helpleness/IMChatRpc/service/rpc/proto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	inithttp()
	//initrpc()

}

func initrpc() {

	config.ConfigInit()
	database.InitClusterClient()
	database.InitMysql()
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	// 创建token验证器
	// 创建认证服务
	authService := auth.NewGRPCAuth()
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor(authService)),
		grpc.StreamInterceptor(middleware.StreamAuthInterceptor(authService)),
	)
	chatServer := controller.NewChatServer()
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	// 创建并注册健康检查服务
	healthService := controller.NewHealthService()
	controller.RegisterHealthService(grpcServer, healthService)

	// 可选：启用 gRPC Reflection 方便调试
	reflection.Register(grpcServer)

	// 启动定时任务
	go chatServer.StartTimer(context.Background())

	// 服务注册到 Consul
	if err := controller.RegisterWithConsul(); err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}

}

func inithttp() {
	config.ConfigInit()
	database.InitClusterClient()

	grpcServer := grpc.NewServer()
	// 创建并注册健康检查服务
	healthService := controller.NewHealthService()
	controller.RegisterHealthService(grpcServer, healthService)

	// 可选：启用 gRPC Reflection 方便调试
	reflection.Register(grpcServer)

	// 服务注册到 Consul
	if err := controller.RegisterWithConsul(); err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	// 1. 初始化 Service 层
	chatService := service.NewChatService()

	// 2. 启动后台任务：消费本节点的消息队列
	// 使用 context 来优雅地关闭 goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go chatService.StartTimer(ctx)
	log.Println("后台消息队列处理器已启动...")

	// 3. 初始化 API 层 (Handler)，注入 Service
	chatHandler := controller.NewChatHandler(chatService)

	// 4. 初始化 Router 层，注入 Handler
	r := router.SetupRouter(chatHandler)

	// 5. 启动 Gin HTTP 服务器
	log.Println("HTTP SSE 服务器启动于 http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}

}
