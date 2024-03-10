package http

import (
	"dev.farukh/copy-close/http/handlers"
	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	handlers.GroupAuthRequests(router.Group("/auth"))
	router.Run()
}
