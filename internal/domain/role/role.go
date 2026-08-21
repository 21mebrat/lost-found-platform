package role

import (
	"time"

	"github.com/google/uuid"
)

type RoleType string

const (
	RoleUser             RoleType = "user"
	RoleInstitutionStaff RoleType = "institution_staff"
	RoleAdmin            RoleType = "admin"
)

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        RoleType  `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
