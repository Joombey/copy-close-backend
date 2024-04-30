package apimodels

import (
	"dev.farukh/copy-close/models/db_models"
)

type RegisterResponse struct {
	UserID    string         `json:"user_id"`
	AuthToken string         `json:"auth_token"`
	AddressID string         `json:"address_id"`
	Role      db_models.Role `json:"role"`
	ImageURL  string         `json:"image_url"`
}

type UserInfoResponse struct {
	UserID    string              `json:"user_id,omitempty"`
	Login     string              `json:"login,omitempty"`
	AuthToken string              `json:"auth_token,omitempty"`
	Name      string              `json:"name"`
	ImageID   string              `json:"user_image"`
	Role      *db_models.Role     `json:"role,omitempty"`
	Address   *db_models.Address  `json:"address"`
	Services  []db_models.Service `json:"services,omitempty"`
}
