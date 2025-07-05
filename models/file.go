package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Size      int64     `json:"size"`                               // in bytes
	NumChunks int       `json:"num_chunks"`                         // total expected chunks
	Hash      string    `json:"hash"`                               // full file SHA-256 (after complete)
	Status    string    `gorm:"default:'incomplete'" json:"status"` // complete/incomplete
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
