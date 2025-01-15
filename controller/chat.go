package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/helpleness/IMChatRpc/database"
	pb "github.com/helpleness/IMChatRpc/service/rpc/proto/proto"
	"github.com/helpleness/IMChatRpc/utilsware/pool/goroutineware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

type Chatserver struct {
	pb.UnimplementedChatServiceServer
	clients     sync.Map      // 使用 sync.Map 替代普通 map
	redisClient *redis.Client // Redis 客户端
	pool        *goroutineware.Pool
}

// NewChatServer 创建并返回一个新的 Chatserver 实例
func NewChatServer() *Chatserver {
	rdb := database.GetRedisClient() // 获取现有的 Redis 客户端
	pool := goroutineware.Default()  // 获取默认的协程池实例

	return &Chatserver{
		clients:     sync.Map{}, // 初始化 sync.Map
		redisClient: rdb,
		pool:        pool,
	}
}

func (s *Chatserver) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	message := req.GetMessage()
	log.Printf("Received message from %s to %s: %s", message.GetUserFrom(), message.GetSendTarget(), message.GetContent())

	targetUserID := message.GetSendTarget()

	targetIP, err := s.redisClient.Get(ctx, targetUserID).Result()
	if err == redis.Nil {
		log.Printf("Key does not exist")
		return nil, err
	} else if err != nil {
		log.Printf("Error getting value: %v", err)
		return nil, err
	}
	// 获取目标用户的流
	targetStream, exists := s.clients.Load(targetUserID) // 使用 Load 获取值
	// 检查目标用户是否在线
	status, err := s.redisClient.HGet(ctx, targetUserID, "status").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("Error checking user %s status: %v", targetUserID, err)
		return nil, err
	}

	// 如果目标用户在线，直接发送消息
	if status == "offline" {
		// 记录消息
		err := s.storeMessageInRedis(ctx, message)
		if err != nil {
			return &pb.SendMessageResponse{
				Success:      false,
				ErrorMessage: "消息存储失败",
			}, nil
		}
		return &pb.SendMessageResponse{
			Success:      true,
			ErrorMessage: "目标用户暂时不在线,等待其上线后推送",
		}, nil
	}
	if status == "online" && !exists {
		err := s.storeMessageInQueue(ctx, message, targetIP)
		if err != nil {
			return &pb.SendMessageResponse{
				Success:      false,
				ErrorMessage: "消息存储失败",
			}, nil
		}
		return &pb.SendMessageResponse{
			Success:      false,
			ErrorMessage: "目标用户不在此节点",
		}, nil
	}

	// 发送消息到目标用户
	if err := targetStream.(pb.ChatService_ReceiveMessageStreamServer).Send(message); err != nil {
		err := s.storeMessageInRedis(ctx, message)
		if err != nil {
			return &pb.SendMessageResponse{
				Success:      false,
				ErrorMessage: "消息存储失败",
			}, nil
		}
		log.Printf("Error sending message to user %s: %v", targetUserID, err)
		return &pb.SendMessageResponse{
			Success:      false,
			ErrorMessage: "发送消息失败",
		}, nil
	}

	return &pb.SendMessageResponse{
		Success: true,
	}, nil
}

// 存储消息到Redis
func (s *Chatserver) storeMessageInRedis(ctx context.Context, message *pb.MyMessage) error {
	key := fmt.Sprintf("messages:%s", message.GetSendTarget())

	// 将消息序列化为字符串或二进制
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	// 使用 HSET 将消息存储到 Redis 哈希
	return s.redisClient.LPush(context.Background(), key, data, 0).Err()
}

// 将无法处理的消息推送到 Redis 队列中
func (s *Chatserver) storeMessageInQueue(ctx context.Context, message *pb.MyMessage, targetIP string) error {

	queueName := fmt.Sprintf("message_queue" + targetIP) // 队列名，可以按需设置

	// 将消息序列化为二进制数据
	data, err := proto.Marshal(message)
	if err != nil {
		log.Printf("Error serializing message: %v", err)
		return err
	}

	// 使用 LPUSH 将消息推送到队列
	err = s.redisClient.LPush(ctx, queueName, data).Err()
	if err != nil {
		log.Printf("Error pushing message to Redis queue: %v", err)
		return err
	}

	log.Printf("Message added to queue for user %s: %s", message.GetSendTarget(), message.GetContent())
	return nil
}

// 定时任务检查消息队列并重新发送消息
func (s *Chatserver) processMessageQueue(ctx context.Context) {
	IP := viper.GetString("service.ip")
	queueName := fmt.Sprintf("message_queue" + IP) // 队列名，可以按需设置
	// 从队列中阻塞地获取消息
	data, err := s.redisClient.BRPop(ctx, 0, queueName).Result()
	if err != nil {
		log.Printf("Error retrieving message from queue: %v", err)
		return
	}

	// 反序列化消息
	var message pb.MyMessage
	err = proto.Unmarshal([]byte(data[2]), &message)
	if err != nil {
		log.Printf("Error deserializing message: %v", err)
		err = s.redisClient.LPush(ctx, queueName, data).Err()
		if err != nil {
			log.Printf("Error pushing message to Redis queue: %v", err)
			return
		}
		return
	}

	targetUserID := message.GetSendTarget()

	// 检查目标用户是否在线
	targetStream, exists := s.clients.Load(targetUserID) // 使用 Load 获取值
	status, err := s.redisClient.HGet(ctx, targetUserID, "status").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("Error checking user %s status: %v", targetUserID, err)
		err := s.storeMessageInQueue(ctx, &message, IP)
		if err != nil {
			log.Printf("Error re-adding message to queue for user %s: %v", targetUserID, err)
		}
		return
	}

	// 如果目标用户在线，发送消息
	if status == "online" && exists {
		// 这里可以直接调用目标用户的流
		err := targetStream.(pb.ChatService_ReceiveMessageStreamServer).Send(&message)
		if err != nil {
			log.Printf("发送失败 %s: %v", targetUserID, err)
		}
		log.Printf("Re-sent message to user %s: %s", targetUserID, message.GetContent())
	} else {
		// 如果目标用户仍然不在线，将消息放回队列中，等待再次处理
		err := s.storeMessageInQueue(ctx, &message, IP)
		if err != nil {
			log.Printf("Error re-adding message to queue for user %s: %v", targetUserID, err)
		}
	}
	return
}
func (s *Chatserver) StartTimer(ctx context.Context) {
	// 从配置中获取IP
	IP := viper.GetString("service.ip")
	queueName := fmt.Sprintf("message_queue_%s", IP)

	defer s.pool.Release()
	// 循环检查队列并处理消息
	for {
		length, err := s.redisClient.LLen(ctx, queueName).Result()
		if err != nil {
			log.Printf("Error checking queue length: %s", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if length > 0 {
			err := s.pool.Submit(func() {
				s.processMessageQueue(ctx)
			})
			if err != nil {
				log.Printf("Error submitting task to ants pool: %s", err)
			}
		} else {
			// 队列为空，等待一段时间后重试
			time.Sleep(1 * time.Second)
		}
	}
}
func (s *Chatserver) ReceiveMessageStream(req *pb.MessageStreamRequest, stream pb.ChatService_ReceiveMessageStreamServer) error {
	IP := viper.GetString("service.ip")
	userID := req.GetUserId()

	// 关联用户ID和流，并将其标记为在线
	s.clients.Store(userID, stream) // 使用 Store 存储流到 sync.Map

	err := s.redisClient.Set(context.Background(), userID, IP, 0).Err()
	if err != nil {
		log.Printf("Error setting user %s status and stream to online: %v", userID, err)
		return err
	}

	// 检查用户是否已经在线
	status, err := s.redisClient.HGet(context.Background(), userID, "status").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("Error checking user %s status: %v", userID, err)
		return err
	}

	// 如果用户状态是 "online"，则表示该用户已经在线
	if status == "online" {
		log.Printf("User %s is already online", userID)
	} else if errors.Is(err, redis.Nil) {
		// 如果没有找到该用户状态（redis.Nil），则表示用户是首次连接
		err = s.redisClient.HSet(context.Background(), userID, "status", "online").Err()
		if err != nil {
			log.Printf("Error setting user %s status and stream to online: %v", userID, err)
			return err
		}
		log.Printf("User %s connected for message stream", userID)
	} else {
		// 如果状态不是 "online"，也表示某种异常，尝试更新为 online
		err := s.redisClient.HSet(context.Background(), userID, "status", "online").Err()
		if err != nil {
			log.Printf("Error setting user %s status and stream to online: %v", userID, err)
			return err
		}
		log.Printf("User %s connected for message stream", userID)
	}

	// 等待客户端断开连接
	<-stream.Context().Done()

	// 移除用户ID和流的关联
	s.clients.Delete(userID) // 使用 Delete 删除流

	// 更新用户状态为 offline，并移除该用户的 stream
	err = s.redisClient.HSet(context.Background(), userID, "status", "offline").Err()
	if err != nil {
		log.Printf("Error updating user %s status to offline: %v", userID, err)
	}

	log.Printf("User %s disconnected from message stream", userID)
	return nil
}
