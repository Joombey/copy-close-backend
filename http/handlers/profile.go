package handlers

import (
	"errors"
	"net/http"

	"dev.farukh/copy-close/models/errs"

	apimodels "dev.farukh/copy-close/models/api_models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GroupProfileHandlers(rg *gin.RouterGroup) {
	rg.POST("/edit-profile", editProfileHandler)
}

func editProfileHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "expected form data request")
		return
	}

	var profileRequest apimodels.EditProfileRequest
	if err = fromString(form.Value["data"][0], &profileRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if valid := (userRepo.CheckTokenValid(profileRequest.UserID, profileRequest.AuthToken)); !valid {
		c.String(http.StatusUnauthorized, "dont have access to edit this profile with such token and id combination")
		return
	}

	fileName := ""
	if len(form.File["image"]) > 0 {
		fileName = uuid.NewV4().String()
		go c.SaveUploadedFile(form.File["image"][0], getPathForJPEG(fileName))
	}

	err = userRepo.EditProfile(profileRequest, fileName)
	if errors.Is(err, errs.ErrInvalidLoginOrAuthToken) {
		c.String(
			http.StatusUnauthorized,
			"dont have access to edit this profile with such token and id combination",
		)
	} else if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
}
