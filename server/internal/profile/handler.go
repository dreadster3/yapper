package profile

import (
	"errors"
	"net/http"

	"github.com/dreadster3/yapper/server/internal/platform/auth"
	"github.com/gin-gonic/gin"
)

type ProfileHandler interface {
	Create(c *gin.Context)
}

type profileHandler struct {
	repository ProfileRepository
}

func NewProfileHandler(repository ProfileRepository) ProfileHandler {
	return &profileHandler{
		repository: repository,
	}
}

func (h *profileHandler) Create(c *gin.Context) {
	var profile *Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()

	profile.UserId = auth.GetUserIdFromContext(c)

	if err := h.repository.Create(ctx, profile); err != nil {
		if errors.Is(err, ErrProfileAlreadyCreated) {
			c.Status(http.StatusConflict)
		}

		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, profile)
}
