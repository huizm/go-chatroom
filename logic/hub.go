package logic

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
			if _, ok := h.Clients[client]; !ok {
				h.Clients[client] = struct{}{}
			}
		case client := <-h.UnregisterCh:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.SendCh)
			}
		case message := <-h.BroadcastCh:
			h.handleMessage(message)
		}
	}
}

func (h *Hub) handleMessage(msg *Message) {

}
