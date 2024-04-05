package db_models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Address struct {
	CopyCloseBase
	AddressName string
	Lat         float32
	Lon         float32

	Users []User
}

func (u *Address) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return nil
}