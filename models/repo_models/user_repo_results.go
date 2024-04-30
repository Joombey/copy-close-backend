package repomodelsgo

import dbModels "dev.farukh/copy-close/models/db_models"

type RegisterResult struct {
	UserID    string        `json:"user_id"`
	AddressID string        `json:"address_id"`
	AuthToken string        `json:"auth_token"`
	UserImage string        `json:"user_image"`
	Role      dbModels.Role `json:"role"`
}

type UserInfoResult struct {
	User     dbModels.User
	Role     dbModels.Role
	Address  dbModels.Address
	Services []dbModels.Service
}