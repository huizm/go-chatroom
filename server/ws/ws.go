package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/huizm/go-chatroom/logic"
	"log"
	"net/http"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failure: ", err)
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		log.Println("Claims not exists")
		return
	}

	client := &logic.Client{
		User:   claims.(*logic.Claims).User,
		Conn:   conn,
		Hub:    logic.Lobby,
		SendCh: make(chan *logic.Message, 10),
	}

	select {
	case logic.Lobby.RegisterCh <- client:
	default:
		log.Println("Lobby register full. Closing connection...")
		_ = conn.Close()
		return
	}

	go client.WritePump()
	go client.ReadPump()
}
