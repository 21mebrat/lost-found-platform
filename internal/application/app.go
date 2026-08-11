package application

import (
	"log/slog"
	"net/http"

	"github.com/21mebrat/lost-found-platform/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config  config.Config
	DB      *pgxpool.Pool
	Redis   *redis.Client
	Logger  *slog.Logger
	Router  http.Handler
	Version string
}

func NewApp(
	cfg config.Config,
	db *pgxpool.Pool,
	redis *redis.Client,
	logger *slog.Logger,
	version string,
) *App {
	app := &App{
		Config:  cfg,
		DB:      db,
		Redis:   redis,
		Logger:  logger,
		Version: version,
	}
	app.LoadRoutes()

	return app
}
