package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	Base
	FirstName  string
	SecondName string
	Password   string
	RoleID     uint32 `gorm:"foreignkey:role_id"`
	Location
}

type Location struct {
	Address string
	Lat     float32
	Lon     float32
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return nil
}
