package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/gwolves/feedy/internal/app/http"
	"github.com/gwolves/feedy/internal/channeltalk"
	"github.com/gwolves/feedy/internal/config"
	"github.com/gwolves/feedy/internal/feed/adapter"
	"github.com/gwolves/feedy/internal/service"
)

func MustInitUsecase() *service.UseCase {
	cfg := config.MustConfig()
	logger := initLogger(cfg)

	return initUsecase(cfg, logger)
}

func MustInitHTTPServer() *http.Server {
	cfg := config.MustConfig()
	logger := initLogger(cfg)

	u := initUsecase(cfg, logger)

	return http.NewServer(cfg.HTTP.Port, u, logger)
}

func initUsecase(cfg *config.Config, logger *slog.Logger) *service.UseCase {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.Postgres.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	repo := adapter.NewPostgresRepo(conn)
	client := channeltalk.NewClient(cfg.AppSecret, logger)

	return service.NewUseCase(cfg.AppName, repo, client, logger)
}

func initLogger(cfg *config.Config) *slog.Logger {
	var level slog.Leveler
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
