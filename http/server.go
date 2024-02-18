package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func Init() {
	router :=gin.Default()
	router.GET("/123", )
	router.Run()

	uuid1 := uuid.NewV4()
	uuid.Equal()
}