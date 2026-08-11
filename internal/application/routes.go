package application

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) LoadRoutes() {
	router := chi.NewRouter()

	// standard production middleares
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(5 * time.Second))

	// define your routes here
	//sample
	router.Get("/health", a.healthHandler)

	a.Router = router
}

func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"status":  "ok",
		"service": app.Config.AppName,
		"env":     app.Config.Environment,
		"version": app.Version,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
