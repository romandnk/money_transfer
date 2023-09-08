package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/romandnk/money_transfer/config"
	zap_logger "github.com/romandnk/money_transfer/internal/logger/zap"
	"github.com/romandnk/money_transfer/internal/server/httpserver"
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

	// initialize http server
	httpServer := httpserver.New(gin.Default(),
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
