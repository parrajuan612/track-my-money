package handlers

import (
	"net/http"
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateMovement(c *gin.Context) {
	// 1. Obtenemos el ID del usuario real desde el Token
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	// 2. Leemos los datos enviados por React
	var req domain.CreateMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// 3. Guardamos enviando el usuario correcto
	if err := h.service.CreateMovement(c.Request.Context(), userUUID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Movimiento creado exitosamente"})
}
