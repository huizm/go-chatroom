package logic

import (
	"github.com/huizm/go-chatroom/model"
	"golang.org/x/net/websocket"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	User   model.User
	Conn   *websocket.Conn
	SendCh chan *Message
}
