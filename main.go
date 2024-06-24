package main

import (
	"ai_bot/cache"
	"ai_bot/handler"
	"ai_bot/middlerware"
	wit_ai "ai_bot/wit-ai"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	godotenv.Load(".env")
	r := gin.Default()

	// 初始化 Redis 客户端
	cache.InitRedis()
	wit_ai.Init(os.Getenv("WIT_AI_TOKEN"))

	//打印环境变量
	fmt.Println("WIT_AI_TOKEN:", os.Getenv("WIT_AI_TOKEN"))
	fmt.Println("REDIS_ADDR:", os.Getenv("REDIS_ADDR"))

	// 注册路由和处理函数
	r.POST("/storePrivateKey", handler.HandleStorePrivateKey)
	r.GET("/ai", middlerware.HasPrivateKeyMiddleware, handler.AiHandler)

	// 启动服务器
	r.Run(":8080")
}
