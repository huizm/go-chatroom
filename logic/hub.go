package logic

import "log"

type Hub struct {
	Clients      map[*Client]struct{}
	BroadcastCh  chan *Message
	RegisterCh   chan *Client
	UnregisterCh chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:      make(map[*Client]struct{}),
		BroadcastCh:  make(chan *Message),
		RegisterCh:   make(chan *Client),
		UnregisterCh: make(chan *Client),
	}
}

var Lobby = NewHub()

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.RegisterCh:
			log.Println("Register client lobby ", client.User.Username)
			if _, ok := h.Clients[client]; !ok {
				h.Clients[client] = struct{}{}
			}
		case client := <-h.UnregisterCh:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.SendCh)
			}
		case message := <-h.BroadcastCh:
			log.Println("Message received lobby ", message.Content)
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
	for client := range h.Clients {
		log.Println("Send message lobby to ", client.User.Username)
		client.SendCh <- msg
	}

	//// broadcast to all by default
	//if msg.Receivers == nil {
	//	for client := range h.Clients {
	//		if client.User.Username != msg.Sender {
	//			msg.Receivers = append(msg.Receivers, client.User.Username)
	//		}
	//	}
	//}
	//
	//for _, client := range msg.Receivers {
	//	select {
	//	case h.Clients[]  <- msg:
	//	default:
	//		close(client.SendCh)
	//		delete(h.Clients, client)
	//	}
	//}
}
