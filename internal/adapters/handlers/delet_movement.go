package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DELETE /api/v1/movements/:id
func (h *Handler) DeleteMovement(c *gin.Context) {
	movID := c.Param("id")

	// Sacamos el ID del usuario del token
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	userUUID, _ := uuid.Parse(userIDStr.(string))

	// Mandamos a eliminar
	if err := h.service.DeleteMovement(c.Request.Context(), movID, userUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el movimiento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movimiento eliminado correctamente"})
}
