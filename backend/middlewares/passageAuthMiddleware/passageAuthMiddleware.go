package passageAuthMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/passageidentity/passage-go"
)

func PassageAuthMiddleware(psg *passage.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := psg.AuthenticateRequest(c.Request)
		if err != nil {
			// 🚨 Authentication failed!
			c.Redirect(302, "/login")
			c.Abort()
			return
		}
		// Authentication successful
		c.Next()
	}
}
