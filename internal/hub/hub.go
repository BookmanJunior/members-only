package hub

import (
	"fmt"
	"sync"
)

type Hub struct {
	sync.Mutex
	Clients    map[*Client]bool
	RegisterCh chan *Client
	unregister chan *Client
	Broadcast  chan WSResponseMessage
	Servers    map[int]map[*Client]bool
}

func CreateNewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		RegisterCh: make(chan *Client),
		unregister: make(chan *Client),
		Broadcast:  make(chan WSResponseMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.RegisterCh:
			h.Register(client)
		case client := <-h.unregister:
			h.Unregister(client)
		case msg := <-h.Broadcast:
			switch {
			case msg.Type == "error":
				h.BroadcastToUser(msg)
			default:
				h.BroadcastToServer(msg)
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	if _, exists := h.Clients[client]; !exists {
		h.Lock()
		h.Clients[client] = true
		h.Unlock()
		fmt.Println("Registered User", client.User.Username)
	}
}

func (h *Hub) Unregister(client *Client) {
	h.Lock()
	delete(h.Clients, client)
	h.Unlock()
	fmt.Println("Unregistered User", client.User.Username)

}

func (h *Hub) BroadcastToUser(msg WSResponseMessage) {
	for client := range h.Clients {
		if msg.Data.User.Id == client.User.Id {
			client.Send <- msg
		}
	}
}

func (h *Hub) BroadcastToServer(msg WSResponseMessage) {
	for client := range h.Clients {
		for i := range client.User.Servers {
			if msg.Data.ServerId == client.User.Servers[i].Id {
				client.Send <- msg
			}
		}
	}
}
