package auth

import (
	"data-storage/src/config"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type AuthMessage struct {
	Token string `json:"token"`
}

func EnsureUserAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Authenticating ws ...")
		conn, ok := ctx.MustGet("conn").(*websocket.Conn)
		if !ok {
			log.Panicln("Failed to get WebSocket connection from context")
			return
		}

		// Listen for the first message to get the token
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			conn.Close()
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Extract the token from the message
		var authMsg AuthMessage
		if err := json.Unmarshal(message, &authMsg); err != nil {
			log.Println("Failed to unmarshal auth message: ", err)
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Invalid authentication message"))
			conn.Close()
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		claims, err := verify(authMsg.Token)
		if err != nil {
			log.Println("Invalid or expired token:", err)
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Invalid or expired token"))
			conn.Close()
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Store claims and token in context to use them
		ctx.Set("claims", claims)
		ctx.Set("token", authMsg.Token)
		ctx.Next()
	}
}

func EnsureBackendAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-Api-Key")

		if apiKey != config.Env.APIKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		ctx.Next()
	}
}
