package user

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Register(*gin.Context)
	Login(*gin.Context)
}

type userHandler struct {
	repository Repository
}

func NewUserHandler(repository Repository) UserHandler {
	return &userHandler{repository}
}

func (handler *userHandler) Register(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (handler *userHandler) Login(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}
