package handlers

import (
	"net/http"

	"errors"

	apimodels "dev.farukh/copy-close/models/api_models"
	"dev.farukh/copy-close/models/errs"
	"github.com/gin-gonic/gin"
)

func GroupInfoHandlers(rg *gin.RouterGroup) {
	rg.GET("/user/:login", getUserInfoHandler)
}

func getUserInfoHandler(c *gin.Context) {
	login := c.Param("login")
	if login == "" {
		c.String(http.StatusBadRequest, "user login must be specified")
		return
	}

	authToken := c.Query("authToken")
	if authToken == "" {
		c.String(http.StatusBadRequest, "auth token must be specified")
		return
	}

	userInfoResult, err := userRepo.GetUser(login, authToken)
	if errors.Is(err, errs.ErrInvalidLoginOrAuthToken) {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(
		http.StatusOK,
		apimodels.UserInfoResponse{
			UserID:    userInfoResult.User.ID.String(),
			Login:     userInfoResult.User.Login,
			AuthToken: userInfoResult.User.AuthToken.String(),
			Name:      userInfoResult.User.FirstName,
			ImageID:   userInfoResult.User.UserImage.String(),
			Role:      userInfoResult.Role,
			Address:   userInfoResult.Address,
		},
	)
}
