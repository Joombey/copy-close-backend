package handlers

import ("net/http"

	"dev.farukh/copy-close/http/websocket"
	apimodels "dev.farukh/copy-close/models/api_models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GroupOrderHandlers(rg *gin.RouterGroup) {
	rg.POST("/create", createOrderHandler)
	rg.GET("/update", manageOrderHandler)
	rg.POST("/report", reportHandler)
	rg.GET("/listen", listenHandler)
}

func listenHandler(c *gin.Context) {
	conn := websocket.ToWebSocket(c)
	worker := &Worker{
		work: func() {
			conn.WriteMessage(1, []byte("1"))
		},
		trigger: make(chan any),
		quit:    make(chan any),
	}
	workerPool.append("2", worker)

	for {
		select {
		case <-worker.quit:
			close(worker.quit)
			close(worker.trigger)
			workerPool.removeWorker("2", worker)
			return
		case <-worker.trigger:
			worker.work()
		}
	}
}

func reportHandler(c *gin.Context) {
	var report Report
	err := c.BindJSON(&report)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	valid := userRepo.CheckTokenValid(report.UserID.String(), report.AuthToken.String())
	if !valid {
		c.String(http.StatusForbidden, "invalid token")
	}

	orderRepo.Report(
		report.OrderID,
		report.UserID,
		report.Message,
	)

	c.Status(http.StatusOK)
	workerPool.trigger("2")
}

func createOrderHandler(c *gin.Context) {
	var request apimodels.OrderCreationRequests
	err := c.BindJSON(&request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	orderRepo.CreateOrder(request)
	workerPool.trigger("2")
}

func manageOrderHandler(c *gin.Context) {
	user, err := userRepo.GetUserInternal(c.Query("user_id"))
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}
	if !*user.Role.CanSell {
		c.String(http.StatusForbidden, "You cannot")
		return
	}

	v, err := stringToInt64(c.Query("state"))
	if err != nil {
		c.String(http.StatusBadRequest, "accpeted should be either true or false")
		return
	}

	if v < 0 || v > 3 {
		c.String(http.StatusBadRequest, "no such state")
		return
	}

	orderRepo.UpdateOrderState(
		uuid.FromStringOrNil(c.Query("order_id")),
		int(v),
	)
	workerPool.trigger("2")
}

type Report struct {
	OrderID   uuid.UUID `json:"order_id"`
	Message   string    `json:"message"`
	UserID    uuid.UUID `json:"user_id"`
	AuthToken uuid.UUID `json:"auth_token"`
}