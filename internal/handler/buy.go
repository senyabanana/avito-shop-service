package handler

import (
	"net/http"

	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) buyItem(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	item := c.Param("item")
	if item == "" {
		entity.NewErrorResponse(c, h.log, http.StatusBadRequest, "item param is required")
		return
	}

	err = h.services.Inventory.BuyItem(c.Request.Context(), userID, item)
	if err != nil {
		entity.NewErrorResponse(c, h.log, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.StatusResponse{
		Status: "item was successfully purchased",
	})
}
