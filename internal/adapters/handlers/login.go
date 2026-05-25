package handlers

import (
	"net/http"
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: revisa el correo y la contraseña"})
		return
	}

	token, user, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Retornamos 401 Unauthorized si la contraseña es incorrecta o el usuario no existe
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login exitoso",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
