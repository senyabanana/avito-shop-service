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

	api := router.Group("/api")
	{
		api.POST("/auth", h.signUp)
		api.GET("/info", h.getInfo)
		api.POST("/sendCoin", h.sendCoin)
		api.GET("/buy/:item", h.buyItem)
	}

	return router
}
