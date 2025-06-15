package profile

import (
	"errors"
	"net/http"

	"github.com/dreadster3/yapper/server/internal/platform/auth"
	"github.com/gin-gonic/gin"
)

const (
	ProfileContextKey = "profile"
)

func GetProfileFromContext(c *gin.Context) *Profile {
	return c.MustGet(ProfileContextKey).(*Profile)
}

func InjectProfileMiddleware(repository ProfileRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := auth.GetUserIdFromContext(c)
		ctx := c.Request.Context()

		profile, err := repository.FindByUserId(ctx, userId)
		if err != nil {
			if errors.Is(err, ErrProfileNotFound) {
				c.Status(http.StatusNotFound)
			}

			c.Error(err)
			c.Abort()
			return
		}

		c.Set(ProfileContextKey, profile)
		c.Next()
	}
}
