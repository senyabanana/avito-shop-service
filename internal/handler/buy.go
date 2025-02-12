package handler

import (
	"net/http"
	
	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) buyItem(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	item := c.Param("item")
	if item == "" {
		entity.NewErrorResponse(c, http.StatusBadRequest, "item param is required")
		return
	}

	err = h.services.Inventory.BuyItem(userID, item)
	if err != nil {
		entity.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
