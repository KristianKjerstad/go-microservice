package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router  http.Handler
	redisDB *redis.Client
	config  Config
}

func New(config Config) *App {
	app := &App{
		config: config,
		redisDB: redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		}),
	}
	app.loadRoutes()
	return app
}

func (app *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.ServerPort),
		Handler: app.router,
	}

	err := app.redisDB.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("Failed to connect to redis: %w", err)
	}

	defer func() {
		if err := app.redisDB.Close(); err != nil {
			fmt.Println("Failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")
	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("Failed to start server %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}

	return nil
}
