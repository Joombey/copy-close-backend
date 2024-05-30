package db_models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Message struct {
	CopyCloseBase
	Text string

	OrderID uuid.UUID
	UserID  uuid.UUID
}

func (c *Message) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewV4()
	return nil
}
