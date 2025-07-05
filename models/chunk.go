package models

import (
	"github.com/google/uuid"
)

type Chunk struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	FileID   uuid.UUID `gorm:"type:uuid;index" json:"file_id"` // foreign key
	Index    int       `json:"index"`                          // 0-based chunk index
	Hash     string    `json:"hash"`                           // SHA-256 of the chunk
	Uploaded bool      `json:"uploaded"`                       // mark upload status
}
