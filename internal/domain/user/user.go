package user

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusBanned    Status = "banned"
	StatusSuspended Status = "suspended"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	MiddleName   string    `json:"middle_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	Email        *string   `json:"email"`
	PasswordHash string    `json:"password_hash"`

	Fayda *string `json:"fayda"`

	LanguageCode string `json:"language_code"`

	PhoneVerified bool `json:"phone_verified"`
	FaydaVerified bool `json:"fayda_verified"`

	ProfileImageURL *string `json:"profile_image_url"`

	Status    Status     `json:"status"`
	Role      Role       `json:"role"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
