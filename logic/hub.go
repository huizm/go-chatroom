package logic

import "log"

type Hub struct {
	Name         string
	Clients      map[uint]*Client
	BroadcastCh  chan *Message
	RegisterCh   chan *Client
	UnregisterCh chan *Client
}

func NewHub(name string) *Hub {
	return &Hub{
		Name:         name,
		Clients:      make(map[uint]*Client),
		BroadcastCh:  make(chan *Message),
		RegisterCh:   make(chan *Client),
		UnregisterCh: make(chan *Client),
	}
}

var Lobby = NewHub("Lobby")

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.RegisterCh:
			log.Printf("[%v (%v)] Register client: %v", h.Name, len(h.Clients), client.User.Username)
			if _, ok := h.Clients[client.User.ID]; !ok {
				h.Clients[client.User.ID] = client
			}
		case client := <-h.UnregisterCh:
			log.Printf("[%v (%v)] Unregister client: %v", h.Name, len(h.Clients), client.User.Username)
			if _, ok := h.Clients[client.User.ID]; ok {
				delete(h.Clients, client.User.ID)
				close(client.SendCh)
			}
		case message := <-h.BroadcastCh:
			log.Printf("[%v (%v)] Receive message from %v: \"%v...\"", h.Name, len(h.Clients), message.Sender.Username, string(message.Content)[:20])
			h.handleMessage(message)
		}
	}
}

func (h *Hub) handleMessage(msg *Message) {
	switch msg.Type {
	case MTText:
		h.handleText(msg)
	case MTInvite:
		h.handleInvite(msg)
	}
}

func (h *Hub) handleInvite(msg *Message) {
	// TODO: implement chatroom creation
}

func (h *Hub) handleText(msg *Message) {
	h.sendMessage(msg)
}

func (h *Hub) sendMessage(msg *Message) {
	for _, c := range h.Clients {
		if c.User != msg.Sender {
			log.Printf("[%v (%v)] Send message to %v: \"%v...\"", h.Name, len(h.Clients), c.User.Username, string(msg.Content)[:20])
			c.send(msg)
		}
	}
}
