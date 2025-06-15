package router

import (
	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/dreadster3/yapper/server/internal/messages"
	"github.com/dreadster3/yapper/server/internal/platform/router/middleware"
	"github.com/dreadster3/yapper/server/internal/profile"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

func SetupRouter(translator ut.Translator, jwtConfig *middleware.JWTConfig, profileRepository profile.ProfileRepository, chatHandler chat.ChatHandler, profileHandler profile.ProfileHandler, messageHandler messages.MessageHandler) (*gin.Engine, error) {
	engine := gin.Default()
	engine.Use(middleware.ErrorMiddleware(translator))

	jwtMiddleware, err := middleware.NewJWTMiddleware(*jwtConfig)
	if err != nil {
		return nil, err
	}

	profileMiddleware := profile.InjectProfileMiddleware(profileRepository)
	v1 := engine.Group("/api/v1", jwtMiddleware.Middleware())
	{
		chatRoutes := v1.Group("/chats", profileMiddleware)
		chatRoutes.POST("", chatHandler.Create)
		chatRoutes.POST("/stream", chatHandler.Stream)

		messagesRoutes := chatRoutes.Group("/:chat_id/messages")
		messagesRoutes.POST("", messageHandler.SendMessage)

		profileRoutes := v1.Group("/profiles")
		profileRoutes.POST("", profileHandler.Create)
	}

	return engine, nil
}
