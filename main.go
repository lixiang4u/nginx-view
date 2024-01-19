package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/nginx-view/handler"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	r := gin.Default()
	//r.Static("/", handler.AppRoot())
	r.GET("/config", handler.NginxConfigJson)
	r.NoRoute(handler.NginxConfigViewer)
	port := handler.CheckNextUsefulPort(handler.ParseArgPort())
	go func() { _ = r.Run(fmt.Sprintf(":%d", port)) }()
	//go handler.OpenUrl(fmt.Sprintf("http://127.0.0.1:%d", port))

	select {
	case _sig := <-sig:
		log.Println(fmt.Sprintf("[stop] %v\n", _sig))
	}
}
