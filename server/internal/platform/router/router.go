package router

import (
	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/dreadster3/yapper/server/internal/platform/router/middleware"
	"github.com/dreadster3/yapper/server/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
)

func SetupRouter(translator ut.Translator, jwtConfig *middleware.JWTConfig, chatHandler chat.ChatHandler, userHandler user.UserHandler) (*gin.Engine, error) {
	binding.EnableDecoderDisallowUnknownFields = true

	engine := gin.Default()
	engine.Use(middleware.ErrorMiddleware(translator))

	jwtMiddleware, err := middleware.NewJWTMiddleware(*jwtConfig)
	if err != nil {
		return nil, err
	}

	v1 := engine.Group("/api/v1", jwtMiddleware.Middleware())
	{
		chatRoutes := v1.Group("/chat")
		chatRoutes.POST("/stream", chatHandler.Stream)
	}

	return engine, nil
}
