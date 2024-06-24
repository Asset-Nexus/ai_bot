package handler

import (
	"ai_bot/cache"
	wit_ai "ai_bot/wit-ai"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PrivateKeyRequest struct {
	WalletAddress string `json:"walletAddress"`
	PrivateKey    string `json:"privateKey"`
}

// 升级 HTTP 连接到 WebSocket 连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store private key" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Private key stored successfully"})
}

func AiHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	defer conn.Close()

	for {
		// 读取 WebSocket 消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Failed to read message: "+err.Error()))
			return
		}

		// 从上下文中获取私钥
		privateKey := c.GetString("privateKey")
		logrus.Printf("Cur privateKey: %s\n", privateKey)

		// 处理业务逻辑
		response, err := wit_ai.ParseMessage(string(message), privateKey)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Failed to process data: "+err.Error()))
			return
		}

		// 发送处理结果
		conn.WriteMessage(websocket.TextMessage, []byte(response))
	}
}
