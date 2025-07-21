package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helpleness/IMChatRpc/database" // 假设的 database 包路径
	"github.com/helpleness/IMChatRpc/model"
	"github.com/helpleness/IMChatRpc/utilsware/pool/goroutineware" // 假设的 goroutine 池包路径
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// ChatService 封装了所有聊天相关的业务逻辑
type ChatService struct {
	// clients 存储每个在线用户的消息通道
	// key 是 userID, value 是 chan *model.MyMessage
	clients     sync.Map
	redisClient *redis.Client
	pool        *goroutineware.Pool
}

// NewChatService 创建一个新的 ChatService 实例
func NewChatService() *ChatService {
	// 初始化依赖
	rdb := database.GetRedisClient()
	p := goroutineware.Default()

	return &ChatService{
		clients:     sync.Map{},
		redisClient: rdb,
		pool:        p,
	}
}

func (s *ChatService) SendMessage(ctx context.Context, req model.SendMessageRequest) error {
	message := req.Message
	log.Printf("Received message from %s to %s: %s", message.UserFrom, message.SendTarget, message.Content)

	targetUserID := message.SendTarget

	targetIP, err := s.redisClient.Get(ctx, "ip"+targetUserID).Result()
	//if err == redis.Nil {
	//	log.Printf("Key does not exist")
	//	return nil, err
	//} else
	if err != nil {
		targetIP = ""
		log.Printf("Error getting value: %v", err)
		//return nil, err
	}
	// 获取目标用户的流
	targetChannel, exists := s.clients.Load(targetUserID) // 使用 Load 获取值
	// 检查目标用户是否在线
	status, err := s.redisClient.HGet(ctx, targetUserID, "status").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("Error checking user %s status: %v", targetUserID, err)
		return err
	}
	if errors.Is(err, redis.Nil) {
		status = "offline"
	}

	// 如果目标用户在线，直接发送消息
	if status == "offline" {
		// 记录消息
		err := s.storeMessageInRedis(ctx, message)
		if err != nil {
			return nil
		}
		return nil
	}
	if status == "online" && !exists {
		err := s.storeMessageInQueue(ctx, message, targetIP)
		if err != nil {
			return nil
		}
		return nil
	}

	// 发送消息到目标用户
	select {
	case targetChannel.(chan *model.MyMessage) <- &message:
	default:
		log.Printf("警告: 用户 %s 的消息通道已满，消息将转为离线存储。", targetUserID)
		return s.storeMessageInRedis(ctx, message)
	}

	return nil
}

// 存储消息到Redis
func (s *ChatService) storeMessageInRedis(ctx context.Context, message model.MyMessage) error {
	key := fmt.Sprintf("messages:%s", message.SendTarget)

	// 将消息序列化为字符串或二进制
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 使用 HSET 将消息存储到 Redis 哈希
	return s.redisClient.LPush(context.Background(), key, data).Err()
}

// 将无法处理的消息推送到 Redis 队列中
func (s *ChatService) storeMessageInQueue(ctx context.Context, message model.MyMessage, targetIP string) error {

	queueName := fmt.Sprintf("message_queue" + targetIP) // 队列名，可以按需设置

	// 将消息序列化为二进制数据
	data, err := json.Marshal(message)
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

	log.Printf("Message added to queue for user %s: %s", message.SendTarget, message.Content)
	return nil
}

// 定时任务检查消息队列并重新发送消息
func (s *ChatService) processMessageQueue(ctx context.Context) {
	IP := viper.GetString("service.ip")
	queueName := fmt.Sprintf("message_queue%s", IP) // 队列名，可以按需设置
	// 从队列中阻塞地获取消息
	data, err := s.redisClient.BRPop(ctx, 0, queueName).Result()
	if err != nil {
		log.Printf("Error retrieving message from queue: %v", err)
		return
	}

	// 反序列化消息
	var message model.MyMessage
	err = json.Unmarshal([]byte(data[1]), &message)
	if err != nil {
		log.Printf("Error deserializing message: %v", err)
		err = s.redisClient.LPush(ctx, queueName, data).Err()
		if err != nil {
			log.Printf("Error pushing message to Redis queue: %v", err)
			return
		}
		return
	}

	targetUserID := message.SendTarget

	// 检查目标用户是否在线
	targetChannel, exists := s.clients.Load(targetUserID) // 使用 Load 获取值
	status, err := s.redisClient.HGet(ctx, targetUserID, "status").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("Error checking user %s status: %v", targetUserID, err)
		err := s.storeMessageInQueue(ctx, message, IP)
		if err != nil {
			log.Printf("Error re-adding message to queue for user %s: %v", targetUserID, err)
		}
		return
	}

	// 如果目标用户在线，发送消息
	if status == "online" && exists {
		select {
		case targetChannel.(chan *model.MyMessage) <- &message:
		default:
			log.Printf("警告: 用户 %s 的消息通道已满，消息将转为离线存储。", targetUserID)
			err := s.storeMessageInRedis(ctx, message)
			if err != nil {
				log.Printf("Error re-adding message to redis for user %s: %v", targetUserID, err)
			}
		}
		log.Printf("Re-sent message to user %s: %s", targetUserID, message.Content)
	} else {
		// 如果目标用户仍然不在线，将消息放回队列中，等待再次处理
		err := s.storeMessageInQueue(ctx, message, IP)
		if err != nil {
			log.Printf("Error re-adding message to queue for user %s: %v", targetUserID, err)
		}
	}
	return
}
func (s *ChatService) StartTimer(ctx context.Context) {
	// 从配置中获取IP
	IP := viper.GetString("service.ip")
	queueName := fmt.Sprintf("message_queue%s", IP)

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
func (s *ChatService) ReceiveMessageStream(userID string) (chan *model.MyMessage, error) {
	IP := viper.GetString("service.ip")

	msgChannel := make(chan *model.MyMessage, 100)
	s.clients.Store(userID, msgChannel)

	err := s.redisClient.Set(context.Background(), "ip"+userID, IP, 0).Err()
	if err != nil {
		log.Printf("Error setting user %s status and stream to online: %v", userID, err)
		return msgChannel, err
	}

	// 检查用户是否已经在线
	status, err := s.redisClient.HGet(context.Background(), userID, "status").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("Error checking user %s status: %v", userID, err)
		return msgChannel, err
	}

	// 如果用户状态是 "online"，则表示该用户已经在线
	if status == "online" {
		log.Printf("User %s is already online", userID)
	} else if errors.Is(err, redis.Nil) {
		// 如果没有找到该用户状态（redis.Nil），则表示用户是首次连接
		err = s.redisClient.HSet(context.Background(), userID, "status", "online").Err()
		if err != nil {
			log.Printf("Error setting user %s status and stream to online: %v", userID, err)
			return msgChannel, err
		}
		log.Printf("User %s connected for message stream", userID)
	} else {
		// 如果状态不是 "online"，也表示某种异常，尝试更新为 online
		err := s.redisClient.HSet(context.Background(), userID, "status", "online").Err()
		if err != nil {
			log.Printf("Error setting user %s status and stream to online: %v", userID, err)
			return msgChannel, err
		}
		log.Printf("User %s connected for message stream", userID)
	}
	return msgChannel, nil
}

// DisconnectUser 处理用户断开连接
func (s *ChatService) DisconnectUser(userID string) {
	if ch, ok := s.clients.Load(userID); ok {
		close(ch.(chan *model.MyMessage))
		s.clients.Delete(userID)
	}

	// 更新 Redis 中的状态为 offline，但保留 IP 地址以便追踪
	err := s.redisClient.HSet(context.Background(), userID, "status", "offline").Err()
	if err != nil {
		log.Printf("错误：更新用户 %s 状态为 offline 失败: %v", userID, err)
	}
	log.Printf("用户 %s 已断开", userID)
}
