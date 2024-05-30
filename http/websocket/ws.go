package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
}

func ToWebSocket(c *gin.Context) *websocket.Conn{
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	return conn
}
