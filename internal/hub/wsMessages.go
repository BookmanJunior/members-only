package hub

import "github.com/bookmanjunior/members-only/internal/models"

type WSMessage struct {
	Message   string  `json:"message,omitempty"`
	MessageId int     `json:"message_id,omitempty"`
	UserId    int     `json:"user_id,omitempty"`
	Headers   Headers `json:"headers"`
}

// example
// WSResponseMessage {
// "type": "error" or "message"
// "message": error message - used specifically for error messages
// "data" embeded struct contains data, usually whatever user posted
// "status": error status - corresponds to http status
// }

type WSResponseMessage struct {
	Type       string         `json:"type"`
	Message    string         `json:"message"`
	Data       models.Message `json:"data,omitempty"`
	Headers    Headers        `json:"headers,omitempty"`
	StatusCode int            `json:"status_code"`
}

type Headers struct {
	Method string `json:"method"`
}
