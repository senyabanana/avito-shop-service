package handler

import (
	"net/http"

	"github.com/senyabanana/avito-shop-service/internal/entity"
	
	"github.com/gin-gonic/gin"
)

func (h *Handler) sendCoin(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	var input entity.SendCoinRequest
	if err := c.BindJSON(&input); err != nil {
		entity.NewErrorResponse(c, http.StatusBadRequest, "invalid request data")
		return
	}

	if err := h.services.Transaction.SendCoin(userID, input.ToUser, input.Amount); err != nil {
		entity.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
