package handlers

import (
	"errors"
	"net/http"

	"dev.farukh/copy-close/di"
	"dev.farukh/copy-close/models/errs"
	api "dev.farukh/copy-close/models/api_models"
	"github.com/gin-gonic/gin"
)

var repo = di.GetComponent().UserRepo

func GroupAuthRequests(rg *gin.RouterGroup) {
	rg.POST("/sign-up", signUpHandler)
	rg.POST("/sign-in", signInHandler)
}

func signUpHandler(c *gin.Context) {
	var request api.SignUpRequest
	err := c.BindJSON(&request)
	if err != nil {
		println(err.Error())
	}

	err = repo.CreateUser(request)
	if errors.Is(err, errs.ErrUserExists) {
		c.JSON(http.StatusNotFound, err.Error())
	} else if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(200, "user successfully created")
	}
}

func signInHandler(c *gin.Context) {
}
