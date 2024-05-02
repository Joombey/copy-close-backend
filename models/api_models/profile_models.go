package apimodels

import "dev.farukh/copy-close/models/db_models"

type EditProfileRequest struct {
	UserID           string              `json:"id"`
	AuthToken        string              `json:"auth_token"`
	Name             string              `json:"name"`
	Services         []db_models.Service `json:"services,omitempty"`
	ServicesToDelete []string            `json:"services_to_delete"`
}