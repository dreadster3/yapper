package main

import (
	"flag"
	"fmt"

	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/dreadster3/yapper/server/internal/platform/providers"
	"github.com/dreadster3/yapper/server/internal/platform/router"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "Port to listen on")
}

func main() {
	flag.Parse()

	providers, err := providers.SetupProviders("http://localhost:11434")
	if err != nil {
		panic(err)
	}

	chatHandler := chat.NewChatHandler(providers)

	engine := router.SetupRouter(chatHandler)
	engine.Run(fmt.Sprintf(":%d", port))
}
