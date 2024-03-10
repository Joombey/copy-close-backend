package db_models

import (
	core "dev.farukh/copy-close/models/core_models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	CopyCloseBase
	FirstName  string
	SecondName string
	Login      string
	Password   string
	RoleID     uint
	core.Address
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return nil
}
