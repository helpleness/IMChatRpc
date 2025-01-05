package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

// 定义一个全局变量
var (
	RedisClient *redis.Client
)

// 结构体方法
func InitClusterClient() *redis.Client {

	masteraddr := viper.GetString("redis.masteraddr")
	password := viper.GetString("redis.password")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     masteraddr, // redis服务ip:port
		Password: password,   // redis的认证密码
		DB:       0,          // 连接的database库
		PoolSize: 100,        // 连接池
	})
	fmt.Printf("Connecting Redis : %v\n", masteraddr)

	// go-redis库v8版本相关命令都需要传递context.Context参数,Background 返回一个非空的Context,它永远不会被取消，没有值，也没有期限。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 验证是否连接到redis服务端
	res, err := RedisClient.Ping(ctx).Result()

	if err != nil {
		fmt.Printf("Connect Failed! Err: %v\n", err)
		panic(err)
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Connect Failed! Err: %v\n", err)
		}
	}()
	// 输出连接成功标识
	fmt.Printf("Connect Successful! \nPing => %v\n", res)
	return RedisClient

}

func GetRedisClient() *redis.Client {
	return RedisClient
}

/*

// 程序执行完毕释放资源
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {

		}
	}(redisClient)

*/
