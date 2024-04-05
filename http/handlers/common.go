package handlers

import (
	"fmt"

	"os"

	"github.com/gin-gonic/gin"
)

func GroupFileRequests(rg *gin.RouterGroup) {
	rg.GET("/getfile/:filename", getFileHandler)

}

func getFileHandler(c *gin.Context) {
	path, _ := os.Getwd()
	path = fmt.Sprintf("%s/files/%s", path, c.Param("filename"))
	c.File(path)
}
