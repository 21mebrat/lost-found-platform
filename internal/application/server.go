package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (a *App) Start(ctx context.Context) error {

	serverPort := fmt.Sprintf(":%d", a.Config.ServerPort)

	srv := http.Server{
		Addr:         serverPort,
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	serverError := make(chan error, 1)
	go func() {
		logInfo := fmt.Sprintf("Server is runing on port %v", serverPort)
		a.Logger.Info(logInfo)

		err := srv.ListenAndServe()
		if err != nil && errors.Is(
			err,
			http.ErrServerClosed,
		) {
			serverError <- err
		}
	}()

	select {
	case err := <-serverError:
		return err

	case <-ctx.Done():
		a.Logger.Info("shutdown signal recevied")
	}

	shutdownctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := srv.Shutdown(shutdownctx)
	if err != nil {
		return fmt.Errorf("Server shutdown failed: %w", err)
	}
	a.Logger.Info("Server shutdown gracefully.")

	return nil
}
