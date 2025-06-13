package main

import (
	"flag"
	"fmt"

	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/dreadster3/yapper/server/internal/platform/providers"
	"github.com/dreadster3/yapper/server/internal/platform/router"
	"github.com/dreadster3/yapper/server/internal/user"
	"go.uber.org/zap"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "Port to listen on")
}

func _main() error {
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	userRepository := user.NewMongoRepository(logger.With(zap.String("repository", "user")))
	userHandler := user.NewUserHandler(userRepository)

	providers, err := providers.SetupProviders("http://localhost:11434")
	if err != nil {
		return err
	}

	chatHandler := chat.NewChatHandler(providers)

	engine := router.SetupRouter(chatHandler, userHandler)
	engine.Run(fmt.Sprintf(":%d", port))

	return nil
}

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}
