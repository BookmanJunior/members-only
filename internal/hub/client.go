package hub

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bookmanjunior/members-only/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	User models.User
	Hub  *Hub
	Conn *websocket.Conn
	Send chan WSResponseMessage
	DB   *models.ServerMessageModel
}

func CreateNewClient(user models.User, conn *websocket.Conn, hub *Hub, db *models.ServerMessageModel) *Client {
	return &Client{
		User: user,
		Hub:  hub,
		Conn: conn,
		DB:   db,
		Send: make(chan WSResponseMessage),
	}
}

func (c *Client) Read() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		_, r, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			c.Conn.Close()
		}

		msg := &WSMessage{}
		reader := bytes.NewReader(r)
		err = json.NewDecoder(reader).Decode(msg)
		if err != nil {
			fmt.Println(err)
			c.Conn.Close()
			break
		}

		msg.UserId = c.User.Id

		switch msg.Headers.Method {
		case "POST":
			HandleWSMessagePost(c, msg)
		case "PATCH":
			HandleWSMessageUpdate(c, msg)
		case "DELETE":
			HandleWsMessageDelete(c, msg)
		}

	}
}

func (c Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg := <-c.Send:

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println(err)
				c.Conn.Close()
				return
			}

			if err := json.NewEncoder(w).Encode(msg); err != nil {
				fmt.Println(err)
				c.Conn.Close()
			}

			if err := w.Close(); err != nil {
				fmt.Println(err)
				c.Conn.Close()
			}
		}
	}
}

func (c *Client) AddServer(server models.Server) {
	c.User.Servers = append(c.User.Servers, server)
}

func (c *Client) RemoveServer(serverID int) {
	for i, server := range c.User.Servers {
		if server.Id == serverID {
			c.User.Servers = append(c.User.Servers[:i], c.User.Servers[i+1:]...)
			return
		}
	}
}
