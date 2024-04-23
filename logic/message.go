package logic

type MessageType uint8

const (
	MTInvite MessageType = iota
	MTText
)

type Message struct {
	Sender    *Client     `json:"sender"`
	Type      MessageType `json:"type"`
	Receivers []*Client   `json:"receivers"`
	Content   []byte      `json:"content"`
}
