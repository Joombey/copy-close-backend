package websocket

import (
	// "fmt"

	// "encoding/json"
	// "math"
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func EventSenderWS(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.String(http.StatusBadRequest, "missing user_id")
		return
	}
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	for {
		conn.ReadMessage()
	}
}
