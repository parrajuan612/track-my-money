package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMe(c *gin.Context) {
	// 1. El portero (middleware) ya validó el token y nos dejó el ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	// Aquí idealmente buscarías al usuario en tu BD por su ID.
	// Pero para salir a beta rápido y que mantenga la sesión,
	// puedes devolverle un usuario básico construido con ese ID:
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":   userIDStr,
			"name": "Usuario", // O el nombre real de tu BD
		},
	})
}
