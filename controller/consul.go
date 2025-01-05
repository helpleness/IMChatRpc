package controller

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"log"
)

func RegisterWithConsul() error {

	ip := viper.GetString("consul.ip")
	port := viper.GetString("consul.port")
	ID := viper.GetString("consul.id")
	Name := viper.GetString("consul.name")
	consulAddress := fmt.Sprintf("%s:%s", ip, port)
	rpcIp := "192.168.15.113"

	serviceAddress := fmt.Sprintf("%s:%s", rpcIp, "50051")
	config := api.DefaultConfig()
	config.Address = consulAddress
	client, err := api.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create Consul client: %v", err)
	}

	// 服务注册
	registration := &api.AgentServiceRegistration{
		ID:      ID,                       // 服务 ID
		Name:    Name,                     // 服务名称
		Address: rpcIp,                    // 服务地址
		Port:    50051,                    // 服务端口
		Tags:    []string{"grpc", "chat"}, // 服务标签
		Check: &api.AgentServiceCheck{ // 健康检查配置
			GRPC:       serviceAddress,
			GRPCUseTLS: false,
			Interval:   "30s",
			Timeout:    "15s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service with Consul: %v", err)
	}

	log.Printf("Service registered with Consul: %s (%s)", Name, serviceAddress)
	return nil

}
