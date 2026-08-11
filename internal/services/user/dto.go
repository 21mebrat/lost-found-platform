package user

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type UpdateRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
}

type UserResponse struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	EmailVerified   bool   `json:"email_verified"`
	PhoneVerified   bool   `json:"phone_verified"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Role            string `json:"role"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	LastLoginAt     string `json:"LastLoginAt"`
}
