package utils

import (
	api_models "dev.farukh/copy-close/models/api_models"
	repo_models "dev.farukh/copy-close/models/repo_models"
)

func MapFromRepoInfoResultToInfoResponse(userInfoResult repo_models.UserInfoResult) api_models.UserInfoResponse {
	return api_models.UserInfoResponse{
		UserID:    userInfoResult.User.ID.String(),
		Login:     userInfoResult.User.Login,
		AuthToken: userInfoResult.User.AuthToken.String(),
		Name:      userInfoResult.User.FirstName,
		ImageID:   userInfoResult.User.UserImage.String(),
		Role:      &userInfoResult.Role,
		Address:   &userInfoResult.Address,
	}
}

func MapFromListRepoInfoResultToListInfoResponse(userInfoResult []repo_models.UserInfoResult) []api_models.UserInfoResponse {
	var userInfoResponseList []api_models.UserInfoResponse
	for _, v := range userInfoResult {
		userInfoResponseList = append(userInfoResponseList, MapFromRepoInfoResultToInfoResponse(v))
	}
	return userInfoResponseList
}

func MapFromRepoInfoResultToInfoResponseSafe(userInfoResult repo_models.UserInfoResult) api_models.UserInfoResponse {
	return api_models.UserInfoResponse{
		UserID:    userInfoResult.User.ID.String(),
		Login:     userInfoResult.User.Login,
		Name:      userInfoResult.User.FirstName,
		ImageID:   userInfoResult.User.UserImage.String(),
		Address:   &userInfoResult.Address,
	}
}

func MapFromListRepoInfoResultToListInfoResponseSafe(userInfoResult []repo_models.UserInfoResult) []api_models.UserInfoResponse {
	var userInfoResponseList []api_models.UserInfoResponse
	for _, v := range userInfoResult {
		userInfoResponseList = append(userInfoResponseList, MapFromRepoInfoResultToInfoResponseSafe(v))
	}
	return userInfoResponseList
}