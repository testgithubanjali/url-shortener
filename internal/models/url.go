package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URL struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	User        User      `gorm:"foreignKey:UserID"`
	OriginalURL string    `gorm:"not null"`
	ShortCode   string    `gorm:"uniqueIndex;not null"`
	ClickCount  int       `gorm:"default:0"`
	ExpiresAt   *time.Time
	CreatedAt   time.Time
}

func (u *URL) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
