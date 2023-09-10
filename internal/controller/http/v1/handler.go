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
	signKey  string
}

func NewHandler(services *service.Services, logger logger.Logger, signKey string) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
		signKey:  signKey,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.loggerMiddleware())
	gin.SetMode(gin.ReleaseMode)
	h.engine = router

	auth := router.Group("/auth")
	{
		newUserRoutes(auth, h.services.User)
	}

	api := router.Group("/api")
	{
		version := api.Group("/v1", h.authorizationMiddleware())
		{
			newAccountRoutes(version, h.services.Account)
		}
	}

	return router
}
