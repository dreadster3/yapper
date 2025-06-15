package auth

import (
	"github.com/gin-gonic/gin"
)

const (
	UserIdContextKey     = "user_id"
	UserClaimsContextKey = "user_claims"
)

type UserId string

func GetUserIdFromContext(c *gin.Context) UserId {
	return UserId(c.GetString(UserIdContextKey))
}
