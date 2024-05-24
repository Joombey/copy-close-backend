package handlers

import (
	// "net/http"

	// "errors"

	// "dev.farukh/copy-close/models/errs"
	// utils "dev.farukh/copy-close/utils"
	"net/http"

	apimodels "dev.farukh/copy-close/models/api_models"
	"github.com/gin-gonic/gin"
)

func GroupOrderHandlers(rg *gin.RouterGroup) {
	rg.POST("/create", createOrderHandler)
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
