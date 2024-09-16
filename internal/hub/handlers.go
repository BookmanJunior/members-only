package hub

import (
	"fmt"
	"net/http"
)

func HandleWSMessagePost(client *Client, msg *WSMessage) {
	var responseMessage WSResponseMessage
	res, err := client.DB.Insert(msg.Message, client.User.Id)
	if err != nil {
		ServerError(&responseMessage)
		client.Hub.Broadcast <- responseMessage
		return
	}

	responseMessage.Type = "success"
	responseMessage.Message = "Successfully posted your message"
	responseMessage.StatusCode = http.StatusOK
	responseMessage.Data = res

	client.Hub.Broadcast <- responseMessage
}

func HandleWSMessageUpdate(client *Client, msg *WSMessage) {
	var responseMessage WSResponseMessage

	if client.User.Id != msg.UserId {
		UnauthorizedError(&responseMessage)
		client.Hub.Broadcast <- responseMessage
		return
	}

	updatedMessage, err := client.DB.UpdateMessage(msg.MessageId, msg.Message)
	if err != nil {
		ServerError(&responseMessage)
		client.Hub.Broadcast <- responseMessage
		return
	}

	responseMessage.Type = "success"
	responseMessage.Message = "Successfully updated message"
	responseMessage.StatusCode = http.StatusOK
	responseMessage.Data = updatedMessage
	client.Hub.Broadcast <- responseMessage
}

func HandleWsMessageDelete(client *Client, msg *WSMessage) {
	var responseMessage WSResponseMessage

	if client.User.Id != msg.UserId {
		UnauthorizedError(&responseMessage)
		client.Hub.Broadcast <- responseMessage
		return
	}

	err := client.DB.Delete(msg.MessageId)
	if err != nil {
		ServerError(&responseMessage)
		client.Hub.Broadcast <- responseMessage
		return
	}

	responseMessage.Type = "success"
	responseMessage.Message = "Successfully deleted message"
	responseMessage.StatusCode = http.StatusOK
	responseMessage.Message = fmt.Sprintf("Deleted message: %v", msg.MessageId)
	client.Hub.Broadcast <- responseMessage
}
