package hub

import (
	"fmt"
	"sync"
)

type Hub struct {
	sync.Mutex
	clients    map[*Client]bool
	RegisterCh chan *Client
	unregister chan *Client
	Broadcast  chan WSResponseMessage
}

func CreateNewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
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
			for client := range h.clients {
				client.Send <- msg
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	if _, exists := h.clients[client]; !exists {
		h.Lock()
		h.clients[client] = true
		h.Unlock()
		fmt.Println("Registered User", client.User.Username)
	}
}

func (h *Hub) Unregister(client *Client) {
	h.Lock()
	delete(h.clients, client)
	h.Unlock()
	fmt.Println("Unregistered User", client.User.Username)

}
