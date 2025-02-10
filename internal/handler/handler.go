package handler

import "github.com/gin-gonic/gin"

type Handler struct{}

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
