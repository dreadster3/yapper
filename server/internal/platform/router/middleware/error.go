package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func ErrorMiddleware(translator ut.Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			status := http.StatusInternalServerError
			var message any
			message = err.Error()
			if err.IsType(gin.ErrorTypeBind) {
				status = http.StatusBadRequest
			}

			var validationError validator.ValidationErrors
			if ok := errors.As(err, &validationError); ok {
				status = http.StatusBadRequest
				translations := validationError.Translate(translator)
				message = translations
			}

			c.JSON(status, gin.H{
				"error": message,
			})

		}
	}
}
