package logic

type MessageType uint8

const (
	MTText MessageType = iota
	MTCreateHub
	MTInvite

	MTCreatedHub
)

type Message struct {
	Sender    *Client           `json:"sender,omitempty"`
	Type      MessageType       `json:"type"`
	Receivers map[uint]struct{} `json:"receivers,omitempty"`
	Content   []byte            `json:"content"`
}

// Message.Content base64 encoded
// MTText: plain text, omit sender and receivers
// MTCreateHub: hub name, omit sender and receivers
