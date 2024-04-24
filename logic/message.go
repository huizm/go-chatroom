package logic

import (
	"github.com/huizm/go-chatroom/model"
)

type MessageType uint8

const (
	MTText MessageType = iota
	MTInvite
)

type Message struct {
	Sender *model.User `json:"sender,omitempty"`
	Type   MessageType `json:"type"`
	//Receivers []string    `json:"receivers,omitempty"`
	Content []byte `json:"content"`
}
