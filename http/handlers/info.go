package handlers

import (
	"net/http"

	"errors"

	"dev.farukh/copy-close/models/errs"
	utils "dev.farukh/copy-close/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GroupInfoHandlers(rg *gin.RouterGroup) {
	rg.GET("/user/:login", getUserInfoHandler)
	rg.GET("/order-list/:userID", getOrderInfoHandler)
}

func getUserInfoHandler(c *gin.Context) {
	login := c.Param("login")
	if login == "" {
		c.String(http.StatusBadRequest, "user login/id must be specified")
		return
	}

	authToken := c.Query("authToken")
	if authToken == "" {
		c.String(http.StatusBadRequest, "auth token must be specified")
		return
	}

	userInfoResult, err := userRepo.GetUser(login, authToken, c.Query("userID"))
	if errors.Is(err, errs.ErrInvalidLoginOrAuthToken) {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.MapFromRepoInfoResultToInfoResponse(userInfoResult))
}

func getOrderInfoHandler(c *gin.Context) {
	authToken := c.Query("auth_token")
	userID := c.Param("userID")
	if !userRepo.CheckTokenValid(userID, authToken) {
		c.String(http.StatusUnauthorized, "not valid user id and token combination")
		return
	}

	c.JSON(http.StatusOK, orderRepo.GetOrderInfo(uuid.FromStringOrNil(userID)))
}