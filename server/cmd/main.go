package main

import (
	"flag"
	"fmt"

	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/dreadster3/yapper/server/internal/platform/router"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "Port to listen on")
}

func main() {
	flag.Parse()

	chatHandler := chat.NewChatHandler()

	engine := router.SetupRouter(chatHandler)
	engine.Run(fmt.Sprintf(":%d", port))
}
