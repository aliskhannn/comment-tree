package main

import (
	"context"
	"errors"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"

	"github.com/aliskhannn/comment-tree/internal/api/handlers/comment"
	"github.com/aliskhannn/comment-tree/internal/api/router"
	"github.com/aliskhannn/comment-tree/internal/api/server"
	"github.com/aliskhannn/comment-tree/internal/config"
	commentrepo "github.com/aliskhannn/comment-tree/internal/repository/comment"
	commentsvc "github.com/aliskhannn/comment-tree/internal/service/comment"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	zlog.Init()
	cfg := config.Must()

	// Connect to PostgreSQL master and slave databases.
	opts := &dbpg.Options{
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}

	slaveDNSs := make([]string, 0, len(cfg.Database.Slaves))

	for _, s := range cfg.Database.Slaves {
		slaveDNSs = append(slaveDNSs, s.DSN())
	}
	zlog.Logger.Info().Msgf("db url: %s", cfg.Database.Master.DSN())
	db, err := dbpg.New(cfg.Database.Master.DSN(), slaveDNSs, opts)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Connect to Redis.
	dbNum, err := strconv.Atoi(cfg.Redis.Database)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to parse redis database")
	}

	zlog.Logger.Info().Msgf("redis config: %s, %s, %d", cfg.Redis.Address, cfg.Redis.Password, dbNum)
	rdb := redis.New(cfg.Redis.Address, cfg.Redis.Password, dbNum)

	if err = rdb.Ping(ctx).Err(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to redis")
	}

	// Initialize comment repository, service and handlers.
	repo := commentrepo.NewRepository(db)
	service := commentsvc.NewService(repo)
	handler := comment.NewHandler(service)

	// Start HTTP server
	r := router.New(handler)
	s := server.New(cfg.Server.HTTPPort, r)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for shutdown signal.
	<-ctx.Done()
	zlog.Logger.Info().Msg("shutdown signal received")

	// Graceful shutdown with timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	zlog.Logger.Info().Msg("shutting down server")
	if err := s.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to shutdown server")
	}
	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		zlog.Logger.Info().Msg("timeout exceeded, forcing shutdown")
	}

	// Close master and slave databases.
	if err := db.Master.Close(); err != nil {
		zlog.Logger.Printf("failed to close master DB: %v", err)
	}
	for i, s := range db.Slaves {
		if err := s.Close(); err != nil {
			zlog.Logger.Printf("failed to close slave DB %d: %v", i, err)
		}
	}
}
