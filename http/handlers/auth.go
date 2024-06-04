package handlers

import (
	"errors"
	"net/http"

	"dev.farukh/copy-close/di"
	api "dev.farukh/copy-close/models/api_models"
	"dev.farukh/copy-close/models/errs"
	"github.com/gin-gonic/gin"
)

var userRepo = di.GetComponent().UserRepo
var fileRepo = di.GetComponent().FileRepo
var orderRepo = di.GetComponent().OrderRepo
var chatRepo = di.GetComponent().ChatRepo
var adminRepo = di.GetComponent().AdminRepo

var devToken = "1"

func init() {
	println(devToken)
}

func GroupAuthRequests(rg *gin.RouterGroup) {
	rg.POST("/register", registerHandler)
	rg.POST("/login", logInHandler)
}

func registerHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, expecting multipart form"})
		return
	}

	var request api.RegisterRequest
	if err := fromString(form.Value["register"][0], &request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := userRepo.RegisterUser(request, c.Query("devKey") == devToken)
	if errors.Is(err, errs.ErrUserExists) {
		c.JSON(http.StatusNotFound, err.Error())
	} else if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		err = c.SaveUploadedFile(form.File["image"][0], getPathForJPEG(result.UserImage))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(200, result)
	}
}

func logInHandler(c *gin.Context) {
	var request api.LogInRequest
	c.BindJSON(&request)
	token, err := userRepo.LogInUser(request)
	if errors.Is(err, errs.ErrInvalidLoginOrPassword) {
		c.JSON(http.StatusUnauthorized, err.Error())
	} else if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.String(http.StatusOK, token.String())
	}
}
