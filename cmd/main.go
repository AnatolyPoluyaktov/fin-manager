package main

import (
	"context"
	"fin-manager/internal/config"
	"fin-manager/internal/repositories"
	"fin-manager/internal/storage/pg"
	"fin-manager/internal/transport/http_server"
	"fin-manager/internal/transport/http_server/middleware"
	"fin-manager/internal/usecase"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoadConfig()
	log := setupLogger(cfg.Env)
	auth_token := cfg.AuthToken
	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug mode is enabled")
	storage := pg.New(ctx, *cfg)
	defer storage.Close()
	if err := storage.Migrate(cfg.MigrationsPath); err != nil {
		log.Error(err.Error())
		return
	}

	category_repo, err := repositories.NewCategoryRepository(storage)
	if err != nil {
		log.Error(err.Error())
		return
	}

	expenses_repo, err := repositories.NewExpensesRepository(storage)
	if err != nil {
		log.Error(err.Error())
	}

	category_usecase := usecase.NewCategoryUseCase(category_repo)
	expenses_usecase := usecase.NewExpenseUseCase(expenses_repo)

	category_handler := http_server.NewCategoryHandler(category_usecase)
	expenses_handler := http_server.NewExpenseHandler(expenses_usecase)
	router := chi.NewRouter()
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(auth_token))
		r.Mount("/categories", http_server.NewCategoryRouter(category_handler))
		r.Mount("/expenses", http_server.NewExpenseRouter(expenses_handler))
	})

	// create info route
	router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Fin manager API"))
		if err != nil {
			return
		}

	})
	server := http_server.NewServer(router, cfg.Server.Address)
	// Канал для ожидания сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Server is running on port 8080")
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	// Ожидание сигнала завершения
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
