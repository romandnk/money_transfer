package main

import (
	"context"
	"github.com/romandnk/money_transfer/config"
	v1 "github.com/romandnk/money_transfer/internal/controller/http/v1"
	zap_logger "github.com/romandnk/money_transfer/internal/logger/zap"
	"github.com/romandnk/money_transfer/internal/server/httpserver"
	"github.com/romandnk/money_transfer/internal/service"
	"github.com/romandnk/money_transfer/internal/storage/postgres"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

const configPath = "config/config.yaml"

func main() {
	// initialize config
	cfg, err := config.New(configPath)
	if err != nil {
		panic("Config error: " + err.Error())
	}

	// initialize zap logger
	log, err := zap_logger.New(
		zap_logger.Level(cfg.ZapLogger.Level),
		zap_logger.Encoding(cfg.ZapLogger.Encoding),
		zap_logger.OutputPaths(cfg.ZapLogger.OutputPath),
		zap_logger.ErrorOutputPaths(cfg.ZapLogger.ErrorOutputPath),
	)

	log.Info("using zap logger")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// initialize postgres storage
	storage, err := postgres.NewPostgresStorage(cfg.PostgresDB.URL,
		postgres.MaxPoolSize(cfg.PostgresDB.MaxPoolSize),
		postgres.ConnTimeout(cfg.PostgresDB.ConnTimeout),
		postgres.ConnAttempts(cfg.PostgresDB.ConnAttempts),
	)
	if err != nil {
		log.Error("error initialize postgres storage: %w", err)
		return
	}
	defer storage.Close()

	log.Info("using postgres storage")

	// initialize services
	services := service.NewServices(storage, cfg.JWT.SignKey)

	// initialize handler
	handler := v1.NewHandler(services, log, cfg.JWT.SignKey)

	// initialize http server
	httpServer := httpserver.New(handler.InitRoutes(),
		httpserver.Port(cfg.HTTPServer.Port),
		httpserver.ReadTimeout(cfg.HTTPServer.ReadTimeout),
		httpserver.WriteTimeout(cfg.HTTPServer.WriteTimeout),
		httpserver.ShutdownTimeout(cfg.HTTPServer.ShutdownTimeout),
	)

	log.Info("starting http server...")
	httpServer.Start()

	select {
	case err = <-httpServer.Notify():
		log.Error("error starting http server", zap.Error(err))
		cancel()
	case <-ctx.Done():
		log.Info("stopping http server...")
		err = httpServer.Stop(ctx)
		if err != nil {
			log.Error("error stopping http server", zap.Error(err))
		}

		log.Info("http server stopped")
	}
}
