package handlers

import (
	"net/http"

	"dev.farukh/copy-close/http/websocket"
	"dev.farukh/copy-close/models/db_models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GroupAdminHandlers(rg *gin.RouterGroup) {
	rg.GET("/blocklist", blocklistHandler)
	rg.PUT("/block", blockHandler)
}

func blockHandler(c *gin.Context) {
	var request BlockRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	if valid := userRepo.CheckTokenValid(request.UserID, request.AuthToken); !valid {
		c.String(http.StatusUnauthorized, "wrong token or uid")
		return
	}

	info, _ := userRepo.GetUserInternal(request.UserID)
	if !*info.Role.CanBan {
		c.String(http.StatusForbidden, "not allowed action")
		return
	}

	if c.Query("block") == "true" {
		userRepo.DeleteUser(request.UserBlockID)
		adminRepo.SetSolution(
			uuid.FromStringOrNil(request.ReportID),
			db_models.REPORT_STATE_BLOCK,
		)
	} else {
		adminRepo.SetSolution(
			uuid.FromStringOrNil(request.ReportID),
			db_models.REPORT_STATE_REJECTED,
		)
	}
	workerPool.trigger("1")
}



func blocklistHandler(c *gin.Context) {
	userID := c.Query("user_id")
	token := c.Query("auth_token")
	if valid := userRepo.CheckTokenValid(userID, token); !valid {
		c.String(http.StatusUnauthorized, "wrong token or uid")
		return
	}

	info, _ := userRepo.GetUserInternal(userID)
	if !*info.Role.CanBan {
		c.String(http.StatusForbidden, "not allowed action")
		return
	}

	if shouldListen := c.Query("listen"); shouldListen == "true" {
		conn := websocket.ToWebSocket(c)
		worker := &Worker{
			work: func() {
				conn.WriteMessage(1, []byte("1"))
			},
			trigger: make(chan any),
			quit:    make(chan any),
		}
		workerPool.append(
			"1",
			worker,
		)
		for {
			select {
			case <-worker.quit:
				close(worker.quit)
				close(worker.trigger)
				workerPool.removeWorker("1", worker)
				return
			case <-worker.trigger:
				worker.work()
			}
		}
	}
	blockList := adminRepo.GetBlocklist()
	println(len(blockList))
	c.JSON(http.StatusOK, blockList)
}

type BlockRequest struct {
	UserID    string `json:"user_id"`
	AuthToken string `json:"auth_token"`

	UserBlockID string `json:"user_block_id"`
	ReportID    string `json:"report_id"`
}
