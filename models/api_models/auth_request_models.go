package apimodels

import (
	core "dev.farukh/copy-close/models/core_models"
)

type RegisterRequest struct {
	Name     string       `json:"name" binding:"required"`
	Address  core.Address `json:"address" binding:"required"`
	Login    string       `json:"login" binding:"required"`
	Password string       `json:"password" binding:"required"`
	IsSeller *bool        `json:"is_seller" binding:"required"`
}

type LogInRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
