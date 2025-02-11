package handler

import (
	"github.com/senyabanana/avito-shop-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.POST("/auth", h.authenticate)

		protected := api.Group("/", h.userIdentity)
		{
			protected.GET("/info", h.getInfo)
			protected.POST("/sendCoin", h.sendCoin)
			protected.GET("/buy/:item", h.buyItem)
		}
	}

	return router
}
