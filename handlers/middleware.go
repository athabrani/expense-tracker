package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDCookie, err := c.Cookie("user_id")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Ambil juga cookie username
		usernameCookie, err := c.Cookie("username")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		userID, _ := strconv.Atoi(userIDCookie)
		
		// Simpan userID DAN username di context
		c.Set("userID", userID)
		c.Set("username", usernameCookie)

		c.Next()
	}
}