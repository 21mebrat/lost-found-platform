package user

import (
	"encoding/json"
	"errors"
	"net/http"

	userservice "github.com/21mebrat/lost-found-platform/internal/services/user"
)

type Handler struct {
	service userservice.UserService
}

func NewHandler(service userservice.UserService) *Handler {
	return &Handler{service: service}
}

// Step 1: Send OTP to Phone (Telegram Style)
func (h *Handler) SendOTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req userservice.SendOTPRequest
	decoded := json.NewDecoder(r.Body)
	decoded.DisallowUnknownFields()

	if err := decoded.Decode(&req); err != nil {
		http.Error(w, userservice.ErrorInvalidIput, http.StatusBadRequest)
		return
	}

	resp, err := h.service.SendOTP(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Step 2: Verify OTP Code (Telegram Style)
func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req userservice.VerifyOTPRequest
	decoded := json.NewDecoder(r.Body)
	decoded.DisallowUnknownFields()

	if err := decoded.Decode(&req); err != nil {
		http.Error(w, userservice.ErrorInvalidIput, http.StatusBadRequest)
		return
	}

	resp, err := h.service.VerifyOTP(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, userservice.ErrorOTPNotFound), errors.Is(err, userservice.ErrorInvalidOTP):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, userservice.ErrorTooManyOTPAttempts):
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Step 3: Complete Profile after OTP verification (Telegram Style)
func (h *Handler) CompleteProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req userservice.CompleteProfileRequest
	decoded := json.NewDecoder(r.Body)
	decoded.DisallowUnknownFields()

	if err := decoded.Decode(&req); err != nil {
		http.Error(w, userservice.ErrorInvalidIput, http.StatusBadRequest)
		return
	}

	resp, err := h.service.CompleteProfile(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, userservice.ErrorPhoneNotVerified):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, userservice.ErrorPhoneAlreadyExists),
			errors.Is(err, userservice.ErrorEmailAlreadyExists),
			errors.Is(err, userservice.ErrorFaydaAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
