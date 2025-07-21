package controller

import (
	"encoding/json"
	"fmt"
	"github.com/helpleness/IMChatRpc/model"
	"github.com/helpleness/IMChatRpc/service"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ChatHandler 封装了 ChatService
type ChatHandler struct {
	Service *service.ChatService
}

// NewChatHandler 创建一个新的 Handler
func NewChatHandler(s *service.ChatService) *ChatHandler {
	return &ChatHandler{Service: s}
}

// SendMessage godoc
// @Summary 发送消息
// @Description 向指定用户发送一条消息
// @Tags Chat
// @Accept  json
// @Produce  json
// @Param message body model.SendMessageRequest true "消息内容"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /send [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req model.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体: " + err.Error()})
		return
	}

	// 调用 Service 层处理消息
	err := h.Service.SendMessage(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "消息发送失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "消息已成功发出"})
}

// SSEStream godoc
// @Summary 创建SSE长链接
// @Description 为指定用户ID创建Server-Sent Events长链接以实时接收消息
// @Tags Chat
// @Produce text/event-stream
// @Param userID path string true "用户ID"
// @Success 200 {string} string "Event-Stream"
// @Router /stream/{userID} [get]
func (h *ChatHandler) SSEStream(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	// 获取 http.Flusher
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "不支持流式传输"})
		return
	}

	// 设置 SSE 必需的头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 1. 用户连接，获取消息通道
	msgChan, _ := h.Service.ReceiveMessageStream(userID)

	// 2. 使用 defer 确保在客户端断开时执行清理工作
	defer h.Service.DisconnectUser(userID)

	// 通知客户端连接成功
	writeSSE(c.Writer, "event", "connected", "连接成功，等待消息...")
	flusher.Flush()

	clientGone := c.Request.Context().Done()

	// 3. 循环监听消息通道和客户端断开事件
	for {
		select {
		case <-clientGone:
			log.Printf("客户端 %s 主动断开连接", userID)
			return
		case msg, ok := <-msgChan:
			if !ok {
				// Channel 已被关闭，说明用户已在别处断开
				log.Printf("用户 %s 的通道已被关闭", userID)
				return
			}

			// 将消息序列化为 JSON
			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Printf("错误：序列化消息失败: %v", err)
				continue
			}

			// 以 SSE 格式发送消息
			writeSSE(c.Writer, "message", string(jsonData), "")
			flusher.Flush()
		}
	}
}

// writeSSE 是一个辅助函数，用于写入标准 SSE 格式的数据
func writeSSE(w io.Writer, event, data, id string) {
	if id != "" {
		fmt.Fprintf(w, "id: %s\n", id)
	}
	if event != "" {
		fmt.Fprintf(w, "event: %s\n", event)
	}
	fmt.Fprintf(w, "data: %s\n\n", data)
}
