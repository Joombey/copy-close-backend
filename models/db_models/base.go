package db_models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type CopyCloseBase struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid" json:"id,omitempty"`
	CreatedAt time.Time `json:"-"`
}

func (c *CopyCloseBase) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewV4()
	return nil
}