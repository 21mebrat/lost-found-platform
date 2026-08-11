package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/21mebrat/lost-found-platform/internal/application"
	"github.com/21mebrat/lost-found-platform/internal/cache"
	"github.com/21mebrat/lost-found-platform/internal/config"
	"github.com/21mebrat/lost-found-platform/internal/database"
	"github.com/21mebrat/lost-found-platform/internal/logger"
)

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGALRM,
	)
	defer stop()
	// version
	version := "v0.1.0"

	// load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	// logger
	logger := logger.New(cfg.Environment)
	// get db constructor
	db, err := database.New(ctx, cfg.DatabaseURL)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	// get redis constructor
	redis, err := cache.New(ctx, cfg.RedisAddr)
	if err != nil {
		panic(err)
	}
	defer redis.Close()

	// call the app

	app := application.NewApp(*cfg, db, redis, logger, version)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}
