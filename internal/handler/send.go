package handler

import (
	"net/http"

	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) sendCoin(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	var input entity.SendCoinRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		entity.NewErrorResponse(c, h.log, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := h.services.Transaction.SendCoin(c.Request.Context(), userID, input.ToUser, input.Amount); err != nil {
		entity.NewErrorResponse(c, h.log, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
