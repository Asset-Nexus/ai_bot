package middlerware

import (
	"ai_bot/cache"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func HasPrivateKeyMiddleware(c *gin.Context) {
	walletAddress := c.Query("walletAddress")

	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing wallet address"})
		c.Abort()
		return
	}

	// 从 Redis 中获取私钥
	privateKey, err := cache.GetPrivateKey(walletAddress)
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Private key not found"})
		c.Abort()
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve private key"})
		c.Abort()
		return
	}

	// 将私钥存储到上下文中
	c.Set("privateKey", privateKey)
}
