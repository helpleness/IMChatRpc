package main

import (
	"context"
	"github.com/helpleness/IMChatRpc/config"
	"github.com/helpleness/IMChatRpc/controller"
	"github.com/helpleness/IMChatRpc/database"
	pb "github.com/helpleness/IMChatRpc/service/rpc/proto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	config.ConfigInit()
	database.InitClusterClient()
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
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
