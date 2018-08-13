package tokenmiddleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewHandler creates a new handler for securing parts of a Gin router with token-based access.
func NewHandler(token string) gin.HandlerFunc {
	tokenHeader := "Token " + token

	return func(c *gin.Context) {
		queryToken, exists := c.GetQuery("access_token")

		if !((exists && queryToken == token) || c.GetHeader("Authorization") == tokenHeader) {
			c.Header("WWW-Authenticate", "Token realm=\"Authorization Required\"")
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}
	}
}
