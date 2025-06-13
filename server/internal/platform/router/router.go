package router

import (
	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/dreadster3/yapper/server/internal/user"
	"github.com/gin-gonic/gin"
)

func SetupRouter(chatHandler chat.ChatHandler, userHandler user.UserHandler) *gin.Engine {
	engine := gin.Default()

	engine.POST("/register", userHandler.Register)
	engine.POST("/login", userHandler.Login)

	v1 := engine.Group("/api/v1")
	{
		chatRoutes := v1.Group("/chat")
		chatRoutes.POST("/stream", chatHandler.Stream)
	}

	return engine
}
