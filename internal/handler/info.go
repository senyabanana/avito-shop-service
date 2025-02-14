package handler

import (
	"net/http"

	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getInfo(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		entity.NewErrorResponse(c, h.log, http.StatusUnauthorized, err.Error())
		return
	}

	info, err := h.services.Transaction.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		entity.NewErrorResponse(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, info)
}
