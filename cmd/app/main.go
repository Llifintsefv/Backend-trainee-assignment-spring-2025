package main

import (
	"Backend-trainee-assignment-spring-2025/config"
	"Backend-trainee-assignment-spring-2025/internal/delivery/handler"
	"Backend-trainee-assignment-spring-2025/internal/middleware"
	postgres "Backend-trainee-assignment-spring-2025/internal/repository/postgresql"
	"Backend-trainee-assignment-spring-2025/internal/router"
	"Backend-trainee-assignment-spring-2025/internal/service"
	"Backend-trainee-assignment-spring-2025/pkg/auth"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	db, err := postgres.NewDB(context.Background(), cfg.DBConnStr)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repoJWT := postgres.NewJWTRepo(db, logger, cfg.SecretKey, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	userRepo := postgres.NewUserRepo(db, logger)

	auth := auth.NewAuth(logger, cfg.SecretKey, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	serviceAuth := service.NewAuthService(repoJWT, logger, auth, userRepo)
	authHandler := handler.NewAuthHandler(serviceAuth, logger)
	middleware := middleware.NewAuthMiddleware(auth, logger)

	pvzRepo := postgres.NewPVZRepo(db, logger)
	servicePvz := service.NewPvzService(pvzRepo, logger)
	pvzHandler := handler.NewPvzHandler(servicePvz, logger)

	app := router.NewApp(authHandler, middleware, pvzHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(cfg.Port); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	fmt.Println("Server is running on port", cfg.Port)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Server forced to shutdown: ", "error", err)
		os.Exit(1)
	}

	log.Println("Server exiting")

}
