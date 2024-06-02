package websockets

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebSocketUpgrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// Store the WebSocket connection in the context
		c.Set("conn", conn)
		c.Next()
	}
}
