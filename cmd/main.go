package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/huizm/go-chatroom/model"
	"github.com/huizm/go-chatroom/server/http"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
)

const shutdownTimeout = 5

func init() {
	if err := model.ShouldLoadDB(); err != nil {
		panic(err)
	}
}

func main() {
	gin.SetMode(gin.DebugMode)

	httpServer := http.New()
	go func() {
		log.Println("Listening HTTP on 8080...")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, netHttp.ErrServerClosed) {
			log.Fatalln("HTTP server error: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalln("HTTP server forced to shutdown: ", err)
	}
	log.Println("Exiting server...")
}
