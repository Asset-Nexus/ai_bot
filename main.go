package main

import (
	"ai_bot/cache"
	"ai_bot/handler"
	"ai_bot/middlerware"
	wit_ai "ai_bot/wit-ai"
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

	// 注册路由和处理函数
	r.POST("/storePrivateKey", handler.HandleStorePrivateKey)
	r.GET("/process", middlerware.HasPrivateKeyMiddleware, handler.HandleProcessRequest)

	// 启动服务器
	r.Run(":8080")
}
