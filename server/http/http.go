package http

import (
	"github.com/gin-gonic/gin"
	"github.com/huizm/go-chatroom/server/ws"
	"net/http"
)

func New() *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: newRouter(),
	}
}

func newRouter() (r *gin.Engine) {
	r = gin.Default()

	r.GET("/ws", ws.Handler)

	r.GET("/users", searchUserByUsername)
	r.POST("/users", createUser)

	return
}
