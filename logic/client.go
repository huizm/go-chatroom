package logic

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/huizm/go-chatroom/model"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	User   *model.User
	Conn   *websocket.Conn
	Hub    *Hub
	SendCh chan *Message
}

func (c *Client) ReadPump() {
	defer func() {
		_ = c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket read pump error: ", err)
			}
			break
		}

		if err = c.handleMessage(message); err != nil {
			log.Println("Handle message error: ", err)
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendCh:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}

			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			if _, err = w.Write(message.Content); err != nil {
				return
			}

			for i := 0; i < len(c.SendCh); i++ {
				if _, err = w.Write([]byte{'\n'}); err != nil {
					return
				}
				if _, err = w.Write((<-c.SendCh).Content); err != nil {
					return
				}
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(data []byte) error {
	msg := &Message{}
	if err := json.Unmarshal(data, msg); err != nil {
		log.Println("Unmarshal message error: ", err)
		return err
	}
	msg.Sender = c.User

	c.Hub.BroadcastCh <- msg
	return nil
}

func (c *Client) send(msg *Message) {
	select {
	case c.SendCh <- msg:
	default:
		close(c.SendCh)
		delete(c.Hub.Clients, c.User.ID)
		delete(Lobby.Clients, c.User.ID)
		return
	}
}
