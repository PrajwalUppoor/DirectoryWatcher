package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// ChatMessage represents a message sent in the chat
type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upgrade connection"})
		return
	}
	defer conn.Close()

	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		// Here you can process the message, save it, or send it to another user.
		// For now, let's echo it back to the client.
		err = conn.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
