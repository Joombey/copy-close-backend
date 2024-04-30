package db_models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	CopyCloseBase
	FirstName string
	Login     string
	Password  string

	RoleID    uint
	UserImage uuid.UUID `gorm:"type:uuid"`
	AddressID uuid.UUID `gorm:"type:uuid"`
	AuthToken uuid.UUID `gorm:"type:uuid"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return nil
}

type Service struct {
	CopyCloseBase
	Title  string     `json:"title"`
	Price  uint       `json:"price"`
	UserID *uuid.UUID `gorm:"type:uuid" json:"-"`
}

func (u *Service) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return nil
}
