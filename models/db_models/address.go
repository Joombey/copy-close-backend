package db_models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Address struct {
	CopyCloseBase
	AddressName string  `json:"address"`
	Lat         float32 `json:"lat"`
	Lon         float32 `json:"lon"`

	Users []User `json:"-"`
}

func (u *Address) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return nil
}
