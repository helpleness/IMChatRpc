package model

// MessageType 定义消息类型
type MessageType int32

const (
	TEXT           MessageType = 0
	IMAGE          MessageType = 1
	FILE           MessageType = 2
	FRIEND_REQUEST MessageType = 3
	GROUP_INVITE   MessageType = 4
	ONLINE_STATUS  MessageType = 5
)

// MyMessage 对应 proto 中的 MyMessage
// 我们使用 json tag 来进行 HTTP body 的绑定和序列化
type MyMessage struct {
	MessageID  string      `json:"message_id"`
	UserFrom   string      `json:"user_from"`
	SendTarget string      `json:"send_target"`
	Content    string      `json:"content"`
	Type       MessageType `json:"type"`
	SendTime   int64       `json:"send_time"`
}

// SendMessageRequest 对应 proto 中的 SendMessageRequest
type SendMessageRequest struct {
	Message MyMessage `json:"message" binding:"required"`
}
