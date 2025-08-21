package paste

import (
	"time"

	"github.com/google/uuid"
)

type Paste struct {
	ID        uuid.UUID  `gorm:"column:id" json:"id"`
	Title     string     `gorm:"column:title" json:"title"`
	Content   string     `gorm:"column:content" json:"content"`
	Lang      string     `gorm:"column:lang" json:"lang"`
	ExpiresAt *time.Time `gorm:"column:expires_at" json:"expires_at,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
}
