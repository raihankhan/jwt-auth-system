package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user model in the database.
type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"` // We will store the hashed password
	Email        string `gorm:"uniqueIndex"`
	FullName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"` // For soft delete if needed
}
