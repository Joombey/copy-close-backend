package handlers

import (
	"net/http"

	ws "dev.farukh/copy-close/http/websocket"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GroupChatHandlers(rg *gin.RouterGroup) {
	rg.GET("/messages/:orderID", getMessageHandler)
	rg.POST("/messages/:orderID", createMessageHandler)
}

func getMessageHandler(c *gin.Context) {
	orderID := uuid.FromStringOrNil(c.Param("orderID"))
	if orderID == uuid.Nil {
		c.String(http.StatusBadRequest, "order ID is required")
		return
	}

	if c.Query("chat") == "true" {
		conn := ws.ToWebSocket(c)
		worker := &Worker{
			work: func() {
				conn.WriteMessage(1, []byte("1"))
			},
			trigger: make(chan any),
			quit:    make(chan any),
		}
		workerPool.append(orderID.String(), worker)

		for {
			select {
			case <-worker.quit:
				close(worker.quit)
				close(worker.trigger)
				workerPool.removeWorker(orderID.String(), worker)
				return
			case <-worker.trigger:
				worker.work()
			}
		}
	}

	c.JSON(http.StatusOK, chatRepo.GetMessagesForOrder(orderID))
}

func createMessageHandler(c *gin.Context) {
	orderID := uuid.FromStringOrNil(c.Param("orderID"))
	if orderID == uuid.Nil {
		c.String(http.StatusBadRequest, "order ID is required")
		return
	}

	var request ChatMessageRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.String(http.StatusBadRequest, "wrong request model")
		return
	}

	if valid := userRepo.CheckTokenValid(request.UserID, request.AuthToken); !valid {
		c.String(http.StatusForbidden, "user ID or auth token is not valid")
		return
	}

	chatRepo.CreateMessage(
		orderID,
		uuid.FromStringOrNil(request.UserID),
		request.Text,
	)
	
	workerPool.trigger(orderID.String())
}
