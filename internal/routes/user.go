package routes

import (
	userhandler "github.com/21mebrat/lost-found-platform/internal/handlers/user"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(
	r chi.Router,
	handler *userhandler.Handler,
) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/otp/send", handler.SendOTP)
		r.Post("/otp/verify", handler.VerifyOTP)
		r.Post("/profile/complete", handler.CompleteProfile)
	})
}
