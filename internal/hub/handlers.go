package hub

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/bookmanjunior/members-only/internal/models"
)

func HandleWSMessagePost(client *Client, msg *WSMessage) {
	var responseMessage WSResponseMessage
	newServerMessage := models.ServerMessage{
		ServerId:  msg.ServerID,
		ChannelId: msg.ChannelId,
		Message:   msg.Message,
		SentDate:  time.Now(),
		User:      client.User,
	}

	newServerMessage.User.Servers = nil

	res, err := client.DB.Insert(newServerMessage)
	if err != nil {
		ServerError(&responseMessage)
		responseMessage.Data.User.Id = client.User.Id
		client.Hub.Broadcast <- responseMessage
		return
	}

	responseMessage.Type = "message"
	responseMessage.StatusCode = http.StatusOK
	responseMessage.Data = res

	client.Hub.Broadcast <- responseMessage
}

func HandleWSMessageUpdate(client *Client, msg *WSMessage) {
	var responseMessage WSResponseMessage
	responseMessage.Data.ServerId = msg.ServerID
	responseMessage.Data.ChannelId = msg.ChannelId
	responseMessage.Data.Id = msg.MessageId
	responseMessage.Data.User = client.User
	responseMessage.Data.User.Servers = nil

	updatedMessage, err := client.DB.Update(msg.Message, msg.MessageId, client.User.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			BadRequestError(&responseMessage)
			client.Hub.Broadcast <- responseMessage
			return
		}
		responseMessage.Data.User.Id = client.User.Id
		ServerError(&responseMessage)
		client.Hub.Broadcast <- responseMessage
		return
	}

	responseMessage.Type = "message"
	responseMessage.StatusCode = http.StatusOK
	responseMessage.Data = updatedMessage
	responseMessage.Headers.Method = "PATCH"
	client.Hub.Broadcast <- responseMessage
}

func HandleWsMessageDelete(client *Client, msg *WSMessage) {
	var responseMessage WSResponseMessage
	responseMessage.Data.ServerId = msg.ServerID
	responseMessage.Data.ChannelId = msg.ChannelId
	responseMessage.Data.Id = msg.MessageId
	responseMessage.Data.User = client.User
	responseMessage.Data.User.Servers = nil

	_, err := client.DB.Delete(msg.MessageId, msg.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			UnauthorizedError(&responseMessage)
			client.Hub.Broadcast <- responseMessage
			return
		}
		ServerError(&responseMessage)
		client.Hub.Broadcast <- responseMessage
		return
	}

	responseMessage.Type = "message"
	responseMessage.StatusCode = http.StatusOK
	responseMessage.Message = fmt.Sprintf("Deleted message: %v", msg.MessageId)
	client.Hub.Broadcast <- responseMessage
}
