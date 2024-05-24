package db_models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Order struct {
	CopyCloseBase
	UserID  uuid.UUID `gorm:"type:uuid"`
	Comment *string

	Services []Service `gorm:"many2many:order_services;"`
}

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
