package user

import (
	userrepo "github.com/21mebrat/lost-found-platform/internal/repository/user"
)

// Step 1: Send OTP
type SendOTPRequest struct {
	Phone string `json:"phone"`
}

type SendOTPResponse struct {
	Message          string `json:"message"`
	Phone            string `json:"phone"`
	IsRegistered     bool   `json:"is_registered"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

// Step 2: Verify OTP
type VerifyOTPRequest struct {
	Phone   string `json:"phone"`
	OTPCode string `json:"otp_code"`
}

type VerifyOTPResponse struct {
	Message           string        `json:"message"`
	Phone             string        `json:"phone"`
	IsRegistered      bool          `json:"is_registered"`
	RegistrationToken string        `json:"registration_token,omitempty"` // Provided if user is NOT yet registered
	User              *UserResponse `json:"user,omitempty"`               // Provided if user IS already registered
}

// Step 3: Complete Profile (Telegram-style: Only after phone is verified via registration_token)
type CompleteProfileRequest struct {
	RegistrationToken string  `json:"registration_token"`
	Phone             string  `json:"phone"`
	FirstName         string  `json:"first_name"`
	MiddleName        string  `json:"middle_name"`
	LastName          string  `json:"last_name"`
	Email             *string `json:"email,omitempty"`
	Fayda             *string `json:"fayda,omitempty"`
	LanguageCode      string  `json:"language_code,omitempty"`
	Password          string  `json:"password"`
}

type UpdateRequest struct {
	FirstName       string  `json:"first_name,omitempty"`
	MiddleName      string  `json:"middle_name,omitempty"`
	LastName        string  `json:"last_name,omitempty"`
	Phone           string  `json:"phone,omitempty"`
	Email           *string `json:"email,omitempty"`
	Fayda           *string `json:"fayda,omitempty"`
	LanguageCode    string  `json:"language_code,omitempty"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
}

type UserResponse struct {
	ID              string  `json:"id"`
	FirstName       string  `json:"first_name"`
	MiddleName      string  `json:"middle_name"`
	LastName        string  `json:"last_name"`
	Phone           string  `json:"phone"`
	Email           *string `json:"email,omitempty"`
	Fayda           *string `json:"fayda,omitempty"`
	LanguageCode    string  `json:"language_code"`
	PhoneVerified   bool    `json:"phone_verified"`
	FaydaVerified   bool    `json:"fayda_verified"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	TotalPages int             `json:"total_pages"`
}

func ToUserListResponse(paginated *userrepo.PaginatedUsers) *UserListResponse {
	if paginated == nil {
		return &UserListResponse{
			Users: []*UserResponse{},
		}
	}

	responses := make([]*UserResponse, 0, len(paginated.Users))
	for _, u := range paginated.Users {
		responses = append(responses, ToUserResponse(u))
	}

	return &UserListResponse{
		Users:      responses,
		TotalCount: paginated.TotalCount,
		Page:       paginated.Page,
		TotalPages: paginated.TotalPages,
	}
}
