package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"money_transfer/config"
	"money_transfer/internal/server/httpserver"
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

	httpServer.Start()

	select {
	case <-httpServer.Notify():
		cancel()
	case <-ctx.Done():
		err = httpServer.Stop(ctx)
		if err != nil {
			return
		}
	}
}
