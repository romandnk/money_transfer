package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/money_transfer/internal/logger"
	"github.com/romandnk/money_transfer/internal/service"
)

type Handler struct {
	engine   *gin.Engine
	services *service.Services
	logger   logger.Logger
}

func NewHandler(services *service.Services, logger logger.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.loggerMiddleware())
	gin.SetMode(gin.ReleaseMode)
	h.engine = router

	api := router.Group("/api")
	{
		version := api.Group("/v1")
		{
			newAccountRoutes(version, h.services.Account)
		}
	}

	return router
}
