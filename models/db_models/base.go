package db_models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type CopyCloseBase struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
}
