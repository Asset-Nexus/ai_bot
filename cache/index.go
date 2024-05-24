package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

// 初始化 Redis 客户端
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // Redis 服务器地址
		Password: os.Getenv("REDIS_PWD"),  // Redis 密码
		DB:       0,                       // 使用默认数据库
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		//log.Error("connect to redis fail", "err", err)
		panic(fmt.Errorf("connect to redis fail, err is %v", err))
	} else {
		fmt.Println("connect to redis success")
	}
}

// 存储私钥
func StorePrivateKey(walletAddress string, privateKey string) error {
	return rdb.Set(ctx, walletAddress, privateKey, 30*time.Minute).Err()
}

// 获取私钥
func GetPrivateKey(walletAddress string) (string, error) {
	return rdb.Get(ctx, walletAddress).Result()
}
