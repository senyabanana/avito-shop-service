package handler

import (
	"net/http"

	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) authenticate(c *gin.Context) {
	var input entity.AuthRequest

	if err := c.BindJSON(&input); err != nil {
		entity.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.services.Authorization.GetUser(input.Username)
	if err != nil {
		if err := h.services.Authorization.CreateUser(input.Username, input.Password); err != nil {
			entity.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		entity.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.AuthResponse{Token: token})
}
