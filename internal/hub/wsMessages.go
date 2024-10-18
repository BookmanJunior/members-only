package hub

import "github.com/bookmanjunior/members-only/internal/models"

// example
// WSResponseMessage {
// "type": "error" or "message"
// "message": error message - used specifically for informational messages
// "data" embeded struct contains data, usually whatever user posted
// "status": error status - corresponds to http status
// }

type WSMessage struct {
	Type      string  `json:"type"`
	ServerID  int     `json:"server_id"`
	ChannelId int     `json:"channel_id"`
	Message   string  `json:"message,omitempty"`
	MessageId int     `json:"message_id,omitempty"`
	UserId    int     `json:"user_id,omitempty"`
	Headers   Headers `json:"headers"`
}

type WSResponseMessage struct {
	Type       string               `json:"type"`
	Message    string               `json:"message,omitempty"`
	Data       models.ServerMessage `json:"data,omitempty"`
	Headers    Headers              `json:"headers,omitempty"`
	StatusCode int                  `json:"status_code"`
}

type Headers struct {
	Method string `json:"method"`
}
