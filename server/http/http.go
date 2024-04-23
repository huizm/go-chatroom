package http

import (
	"github.com/gin-gonic/gin"
	"github.com/huizm/go-chatroom/server/middleware"
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

	r.GET("/users", searchUserByUsername) // ?username=cohad
	r.POST("/users", createUser)          // {"username": "cohad", "password": "hashed"}
	r.POST("/auth", auth)                 // {"username": "cohad", "password": "hashed"}

	r.Use(middleware.JWT()) // (in header) "Authorization": token
	r.GET("/ws", ws.Handler)

	return
}
