package handler

import (
	"ai_bot/cache"
	wit_ai "ai_bot/wit-ai"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PrivateKeyRequest struct {
	WalletAddress string `json:"walletAddress"`
	PrivateKey    string `json:"privateKey"`
}

// 处理存储私钥的请求
func HandleStorePrivateKey(c *gin.Context) {
	var req PrivateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 将私钥存储到 Redis
	if err := cache.StorePrivateKey(req.WalletAddress, req.PrivateKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store private key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Private key stored successfully"})
}

// 处理业务逻辑请求
func HandleProcessRequest(c *gin.Context) {
	command := c.Query("command")

	if command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing command"})
		return
	}

	// 从上下文中获取私钥
	privateKey, _ := c.Get("privateKey")

	// 处理业务逻辑
	err := wit_ai.ParseMessage(command, privateKey.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process data" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data processed successfully"})
}
