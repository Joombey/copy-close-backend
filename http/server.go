package http

import (
	"net/http"
	"dev.farukh/copy-close/di"
	"dev.farukh/copy-close/repos"
	"dev.farukh/copy-close/http/handlers"
	"github.com/gin-gonic/gin"
)

var repo = di.GetComponent().UserRepo.(*repos.UserRepoImpl)

func Init() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	handlers.GroupAuthRequests(router.Group("/auth"))
	handlers.GroupFileRequests(router.Group("/file"))
	handlers.GroupInfoHandlers(router.Group("/info"))
	handlers.GroupMapHandlers(router.Group("/map"))
	handlers.GroupProfileHandlers(router.Group("/profile"))
	handlers.GroupOrderHandlers(router.Group("/order"))
	
	router.GET("/clear", func(ctx *gin.Context) {
		repo.ClearAll()
		ctx.JSON(http.StatusOK, "zaebic")
	})
	
	router.Run()
}
