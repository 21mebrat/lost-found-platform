package application

import (
	"time"

	userhandler "github.com/21mebrat/lost-found-platform/internal/handlers/user"
	userrepo "github.com/21mebrat/lost-found-platform/internal/repository/user"
	"github.com/21mebrat/lost-found-platform/internal/routes"
	userservice "github.com/21mebrat/lost-found-platform/internal/services/user"
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

	// initate user rep
	userRepository := userrepo.NewPostgresRepository(a.DB)
	UserService := userservice.NewService(userRepository)
	userHandler := userhandler.NewHandler(UserService)

	// user routes
	routes.UserRoutes(router, userHandler)

	a.Router = router
}
