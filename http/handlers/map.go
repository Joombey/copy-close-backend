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
	sellers := utils.MapFromListRepoInfoResultToListInfoResponseSafe(userRepo.GetSellers())
	c.JSON(http.StatusOK, sellers)
}