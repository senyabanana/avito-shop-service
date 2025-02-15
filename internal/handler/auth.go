package handler

import (
	"errors"
	"net/http"

	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) authenticate(c *gin.Context) {
	var input entity.AuthRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		entity.NewErrorResponse(c, h.log, http.StatusBadRequest, "invalid request format")
		return
	}

	_, err := h.services.Authorization.GetUser(c.Request.Context(), input.Username)
	if err != nil && errors.Is(err, entity.ErrUserNotFound) {
		if err := h.services.Authorization.CreateUser(c.Request.Context(), input.Username, input.Password); err != nil {
			entity.NewErrorResponse(c, h.log, http.StatusBadRequest, err.Error())
			return
		}
	}

	token, err := h.services.Authorization.GenerateToken(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		entity.NewErrorResponse(c, h.log, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.AuthResponse{Token: token})
}
