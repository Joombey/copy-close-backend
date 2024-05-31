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
	Deleted   bool

	Messages []Message `json:"-"`
	Orders   []Order   `json:"-"`
	Services []Service `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return nil
}

type Service struct {
	CopyCloseBase
	Title   string     `json:"title"`
	Price   uint       `json:"price"`
	Deleted bool       `json:"-"`
	UserID  *uuid.UUID `gorm:"type:uuid" json:"-"`

	Orders []Order `gorm:"many2many:order_services;"`
}

func (u *Service) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return nil
}
