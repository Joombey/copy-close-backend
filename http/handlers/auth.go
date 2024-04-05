package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"encoding/json"

	"dev.farukh/copy-close/di"
	api "dev.farukh/copy-close/models/api_models"
	"dev.farukh/copy-close/models/errs"
	"github.com/gin-gonic/gin"
)

var userRepo = di.GetComponent().UserRepo

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

	result, err := userRepo.RegisterUser(request)
	if errors.Is(err, errs.ErrUserExists) {
		c.JSON(http.StatusNotFound, err.Error())
	} else if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		err = c.SaveUploadedFile(form.File["image"][0], fileDestination(result.UserImage))
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
	err := userRepo.LogInUser(request)
	if errors.Is(err, errs.ErrInvalidLoginOrPassword) {
		c.JSON(http.StatusUnauthorized, err.Error())
	} else if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, "authorization is successful")
	}
}

func fromString(value string, receiver any) error {
	err := json.Unmarshal([]byte(value), receiver)
	if err != nil {
		return err
	}
	return nil
}

func fileDestination(path string) string {
	absolutePath, _ := os.Getwd()
	return fmt.Sprintf("%s/files/%s.jpeg", absolutePath, path)
}
