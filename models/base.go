package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
}
