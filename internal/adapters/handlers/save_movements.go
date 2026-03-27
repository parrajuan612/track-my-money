package handlers

import (
	"net/http"
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SaveMovements(c *gin.Context) {
	var req domain.SaveStatementRequest

	// Gin llena el modelo de dominio directamente
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pasamos el objeto al servicio
	err := h.service.SaveMovements(c.Request.Context(), req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al guardar los movimientos")
		return
	}

	c.String(http.StatusOK, "Movimientos guardados exitosamente")
}
