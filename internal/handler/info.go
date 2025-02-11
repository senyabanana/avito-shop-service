package handler

import (
	"net/http"
	
	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getInfo(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	info, err := h.services.UserTransaction.GetUserInfo(userID)
	if err != nil {
		entity.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, info)
}
