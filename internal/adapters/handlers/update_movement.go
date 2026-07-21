package handlers

import (
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) UpdateMovement(c *gin.Context) {
	// 1. Obtener el ID del movimiento desde la URL
	movIDStr := c.Param("id")
	movementUUID, err := uuid.Parse(movIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID de movimiento inválido"})
		return
	}

	// 2. Leer los datos enviados por React
	var req domain.UpdateMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// 3. Ejecutar la actualización (Pasamos uuid.Nil temporalmente para el usuario)
	if err := h.service.UpdateMovement(c.Request.Context(), movementUUID, uuid.Nil, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}) // Ahora sí veremos si falla!
		return
	}

	c.JSON(200, gin.H{"message": "Movimiento actualizado exitosamente"})
}
