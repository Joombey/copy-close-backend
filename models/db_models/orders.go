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

	Reports  []Report
	Douments []Document
	Messages []Message
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

type Report struct {
	CopyCloseBase
	OrderID   uuid.UUID `gorm:"type:uuid"`
	Solution int       `gorm:"default:0"`
	Message   string    `gorm:"type:text"`
}

func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewV4()
	return nil
}

var (
	REPORT_STATE_IDLE     = 1
	REPORT_STATE_BLOCK    = 2
	REPORT_STATE_REJECTED = 3
)
