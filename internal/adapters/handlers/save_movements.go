package handlers

import (
	"net/http"
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) SaveMovements(c *gin.Context) {
	// 1. Extraemos el usuario real desde el token JWT
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

	// 2. Leemos el JSON (que ahora trae account_id y bank_id)
	var req domain.SaveStatementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Convertimos el AccountID a UUID
	accountUUID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cuenta inválido"})
		return
	}

	// 4. Enviamos todo al servicio
	err = h.service.SaveMovements(c.Request.Context(), userUUID, accountUUID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el extracto: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Extracto y movimientos guardados exitosamente"})
}
