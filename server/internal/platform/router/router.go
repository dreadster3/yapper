package router

import (
	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/gin-gonic/gin"
)

func SetupRouter(chatHandler chat.ChatHandler) *gin.Engine {
	engine := gin.Default()

	v1 := engine.Group("/api/v1")
	{
		chatRoutes := v1.Group("/chat")
		chatRoutes.POST("/stream", chatHandler.Stream)
	}

	return engine
}
