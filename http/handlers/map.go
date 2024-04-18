package handlers

import (
	"net/http"

	utils "dev.farukh/copy-close/utils"
	"github.com/gin-gonic/gin"
)

func GroupMapHandlers(rg *gin.RouterGroup) {
	rg.GET("/sellers", getSellersHandler)
}

func getSellersHandler(c *gin.Context) {
	userID := c.Query("userID")
	authToken := c.Query("authToken")
	if userID == "" || authToken == "" {
		c.String(http.StatusBadRequest, "No userID nor auth token has been specified")
		return
	}

	if valid := userRepo.CheckTokenValid(userID, authToken); valid {
		c.String(http.StatusUnauthorized, "invalid userID or auth token")
		return
	}
	
	sellers := utils.MapFromListRepoInfoResultToListInfoResponseSafe(userRepo.GetSellers())
	c.JSON(http.StatusOK, sellers)
}