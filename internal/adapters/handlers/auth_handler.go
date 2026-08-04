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

	// 2. Buscamos al usuario por su ID en la base de datos
	// (Asegúrate de tener GetUserByID en tu servicio)
	user, err := h.service.GetUserByID(c.Request.Context(), userIDStr.(string))
	
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado en la base de datos"})
		return
	}

	// 3. Devolvemos los datos REALES extraídos de la base de datos
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,  // ¡Aquí viaja el nombre real!
			"email": user.Email,
		},
	})
}
