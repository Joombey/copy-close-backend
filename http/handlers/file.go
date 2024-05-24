package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GroupFileRequests(rg *gin.RouterGroup) {
	rg.GET("/image/:imageID", getFileHandler)
	rg.GET("/create-session", createSessionHandler)
	rg.POST("/upload/:sessionID", uploadHandler)
	rg.GET("/get-document/:documentID", getDocumentHandler)
}

func getDocumentHandler(c *gin.Context) {
	documentID := uuid.FromStringOrNil(c.Param("documentID"))
	if documentID == uuid.Nil {
		c.String(http.StatusBadRequest, "documentID is required")
		return
	}

	offset, err := stringToInt64(c.GetHeader("Range"))
	if err != nil {
		c.String(http.StatusBadRequest, "Range header is required")
		return
	}

	file, err := fileRepo.GetDocument(documentID, offset)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusOK)
	chunkedFile(file, func(chunk []byte, breakFunc func()) {
		c.Writer.Write(chunk)
		c.Writer.Flush()
	})
}

func createSessionHandler(c *gin.Context) {
	path := getPathForName(c.Query("filename"))
	size, _ := stringToInt64(c.Query("length"))

	sessionID := fileRepo.CreateSession(path, size).String()
	c.String(http.StatusOK, sessionID)
}

func uploadHandler(c *gin.Context) {
	sessionID := uuid.FromStringOrNil(c.Param("sessionID"))
	println(sessionID.String())
	if sessionID == uuid.Nil {
		c.String(http.StatusBadRequest, "sessionID id must be specified")
		return
	}

	header, e := c.FormFile("data")
	if e != nil {
		c.String(http.StatusBadRequest, e.Error())
	}

	file, e := header.Open()
	if e != nil {
		c.String(http.StatusBadRequest, e.Error())
	}

	var (
		id  uuid.UUID
		err error
	)
	
	chunkedReader(
		file,
		func(chunk []byte, breakFunc func()) {
			id, err = fileRepo.WriteToSession(sessionID, chunk)
			if err != nil {
				println(err.Error())
				c.String(http.StatusInternalServerError, err.Error())
				breakFunc()
			} else if id != uuid.Nil {
				c.String(http.StatusOK, id.String())
				breakFunc()
			}
		},
	)
}

func getFileHandler(c *gin.Context) {
	path, _ := os.Getwd()
	path = fmt.Sprintf("%s/files/%s.jpeg", path, c.Param("imageID"))
	c.File(path)
}