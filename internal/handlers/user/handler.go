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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req userservice.RegisterRequest

	decoded := json.NewDecoder(r.Body)
	decoded.DisallowUnknownFields()

	err := decoded.Decode(&req)
	if err != nil {
		http.Error(
			w,
			userservice.ErrorInvalidIput,
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.service.Register(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(
			err,
			userservice.ErrorEmailAlreadyExists,
		):
			http.Error(w, userservice.ErrorEmailAlreadyExists.Error(), http.StatusConflict)
		case errors.Is(
			err,
			userservice.ErrorPhoneAlreadyExists,
		):
			http.Error(w, userservice.ErrorPhoneAlreadyExists.Error(), http.StatusConflict)
		case errors.Is(err, userservice.ErrorUserNotFound):
			http.Error(w, userservice.ErrorUserNotFound.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&user)

}
