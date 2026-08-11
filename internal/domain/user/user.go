package user

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive    Status = "ACTIVE"
	StatusSuspended Status = "SUSPENDED"
	StatusBlocked   Status = "BLOCKED"
)

type Role string

const (
	RoleUser      Role = "USER"
	RoleModerator Role = "MODERATOR"
	RoleAdmin     Role = "ADMIN"
)

type User struct {
	ID uuid.UUID

	// Identity
	FirstName string
	LastName  string

	// Authentication
	Email        string
	Phone        string
	PasswordHash string

	// Verification
	EmailVerified bool
	PhoneVerified bool

	// Profile
	ProfileImageURL string

	// Authorization
	Role   Role
	Status Status

	// Authentication Activity
	LastLoginAt *time.Time

	// Audit
	CreatedAt time.Time
	UpdatedAt time.Time
}
