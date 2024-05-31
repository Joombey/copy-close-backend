package handlers

import (
	// "net/http"

	// "errors"

	// "dev.farukh/copy-close/models/errs"
	// utils "dev.farukh/copy-close/utils"
	"net/http"

	apimodels "dev.farukh/copy-close/models/api_models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GroupOrderHandlers(rg *gin.RouterGroup) {
	rg.POST("/create", createOrderHandler)
	rg.GET("/update", manageOrderHandler)
	rg.POST("/report", reportHandler)
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
}

func createOrderHandler(c *gin.Context) {
	var request apimodels.OrderCreationRequests
	err := c.BindJSON(&request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	orderRepo.CreateOrder(request)
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
}

type Report struct {
	OrderID   uuid.UUID `json:"order_id"`
	Message   string    `json:"message"`
	UserID    uuid.UUID `json:"user_id"`
	AuthToken uuid.UUID `json:"auth_token"`
}