package domain

import (
	"time"
)

type User struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash *string   `json:"-"`
	AuthProvider string    `json:"auth_provider"`
	ExternalID   string    `json:"external_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}
