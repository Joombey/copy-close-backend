package db_models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Order struct {
	CopyCloseBase
	UserID  uuid.UUID `gorm:"type:uuid"`
	State   int       `gorm:"default:0"`
	Comment *string

	Douments []Document
	Services []Service `gorm:"many2many:order_services;"`
}

var (
	STATE_REQUESTED = 0
	STATE_ACCEPTED  = 1
	STATE_REJECTED  = 2
	STATE_COMPLETED = 3
)

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.NewV4()
	return nil
}

type Document struct {
	CopyCloseBase
	Name    string
	Path    string
	OrderID *uuid.UUID `gorm:"type:uuid"`
}

type OrderService struct {
	OrderID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	ServiceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Price     uint      `gorm:"type:uint"`
	Amount    uint      `gorm:"type:uint"`
	Title     string    `gorm:"type:text"`
}
