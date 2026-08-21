package application

import (
	"time"

	"github.com/21mebrat/lost-found-platform/internal/config"
	userhandler "github.com/21mebrat/lost-found-platform/internal/handlers/user"
	"github.com/21mebrat/lost-found-platform/internal/notification"
	otprepo "github.com/21mebrat/lost-found-platform/internal/repository/otp"
	userrepo "github.com/21mebrat/lost-found-platform/internal/repository/user"
	"github.com/21mebrat/lost-found-platform/internal/routes"
	smsservice "github.com/21mebrat/lost-found-platform/internal/services/sms"
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

	// initiate otp rep
	otpRepository := otprepo.NewRedisRepository(a.Redis)

	// build notification service
	smsProvider := ""
	if a.Config.AfroMessageToken != "" && a.Config.AfroMessageSenderID != "" && a.Config.AfroMessageURL != "" {
		smsProvider = "africas_talking"
	}

	notifConfig := config.NotificationConfig{
		Timeout: 10 * time.Second,
		SMS: config.SMSConfig{
			Provider: smsProvider,
			APIKey:   a.Config.AfroMessageToken,
			Username: "default",
			SenderID: a.Config.AfroMessageSenderID,
			BaseURL:  a.Config.AfroMessageURL,
		},
	}

	factory := notification.NewFactory(notifConfig)
	providers, err := factory.Build()
	if err != nil {
		panic("failed to build notification providers: " + err.Error())
	}

	notificationService := notification.NewService(providers)

	// initiate sms service
	smsSvc := smsservice.NewService(notificationService)

	UserService := userservice.NewService(userRepository, otpRepository, smsSvc)
	userHandler := userhandler.NewHandler(UserService)

	// user routes
	routes.UserRoutes(router, userHandler)

	a.Router = router
}
