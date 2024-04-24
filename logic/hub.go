package logic

import (
	"log"
	"math/rand"
	"strconv"
)

type Hub struct {
	Name         string
	Clients      map[uint]*Client
	BroadcastCh  chan *Message
	RegisterCh   chan *Client
	UnregisterCh chan *Client
}

func NewHub(name string) *Hub {
	var postfix string
	if name != "Lobby" {
		postfix = "#" + strconv.Itoa(rand.Intn(10000))
		for _, ok := Hubs[name+postfix]; ok; {
			postfix = "#" + strconv.Itoa(rand.Intn(10000))
		}
	} else {
		postfix = "#0000"
	}

	hub := &Hub{
		Name:         name + postfix,
		Clients:      make(map[uint]*Client),
		BroadcastCh:  make(chan *Message),
		RegisterCh:   make(chan *Client),
		UnregisterCh: make(chan *Client),
	}

	Hubs[hub.Name] = hub
	go hub.Run()
	return hub
}

var Lobby = NewHub("Lobby")
var Hubs = map[string]*Hub{}

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
			log.Printf("[%v (%v)] Receive message from %v: \"%v...\"", h.Name, len(h.Clients), message.Sender.User.Username, string(message.Content)[:min(len(string(message.Content)), 20)])
			h.handleMessage(message)
		}
	}
}

func (h *Hub) handleMessage(msg *Message) {
	switch msg.Type {
	case MTText:
		h.handleText(msg)
	case MTCreateHub:
		h.handleCreateHub(msg)
	case MTInvite:
		h.handleInvite(msg)
	}
}

func (h *Hub) handleText(msg *Message) {
	h.sendMessage(msg)
}

func (h *Hub) handleCreateHub(msg *Message) {
	name := string(msg.Content)
	hub := NewHub(name)
	hub.RegisterCh <- msg.Sender

	msg.Sender.SendCh <- &Message{
		Type:    MTCreatedHub,
		Content: []byte(hub.Name),
	}
}

func (h *Hub) handleInvite(msg *Message) {
	// TODO: implement chatroom creation
}

func (h *Hub) sendMessage(msg *Message) {
	for _, c := range h.Clients {
		if c != msg.Sender {
			log.Printf("[%v (%v)] Send message to %v: \"%v...\"", h.Name, len(h.Clients), c.User.Username, string(msg.Content)[:min(len(string(msg.Content)), 20)])
			c.send(msg)
		}
	}
}
