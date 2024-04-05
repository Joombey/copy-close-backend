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
