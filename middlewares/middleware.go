package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"kysuit.net/go-api/models"
)

func DbMiddleWare(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Bearer ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authenticated.",
			})
			c.Abort()
			return
		}

		token := split[1]
		fmt.Printf("Bearer (%v) \n", token)
		isValid, userId := models.IsTokenValid(token)
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authenticated.",
			})
			c.Abort()
		} else {
			c.Set("user_id", userId)
			c.Next()
		}
	}
}
