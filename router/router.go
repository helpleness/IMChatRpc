package router

import (
	"github.com/gin-gonic/gin"
	"github.com/helpleness/IMChatAdmin/middleware"
	"github.com/helpleness/IMChatRpc/controller"
)

// SetupRouter 配置应用的路由
func SetupRouter(chatHandler *controller.ChatHandler) *gin.Engine {
	r := gin.Default()

	// 路由组
	apiGroup := r.Group("/api/v1")
	{
		// 发送消息的端点
		apiGroup.POST("/send", middleware.AuthMiddleWare(), chatHandler.SendMessage)
		// 建立 SSE 长链接的端点
		apiGroup.GET("/stream/:userID", middleware.AuthMiddleWare(), chatHandler.SSEStream)
	}

	return r
}
