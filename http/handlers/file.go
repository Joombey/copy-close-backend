package handlers

import (
	"fmt"

	"os"

	"github.com/gin-gonic/gin"
)

func GroupFileRequests(rg *gin.RouterGroup) {
	rg.GET("/image/:imageID", getFileHandler)
}

func getFileHandler(c *gin.Context) {
	path, _ := os.Getwd()
	path = fmt.Sprintf("%s/files/%s.jpeg", path, c.Param("imageID"))
	c.File(path)
}
