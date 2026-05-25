package handlers

import (
	"net/http"
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LoginGoogle(c *gin.Context) {
	var req domain.GoogleLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere el token de Google"})
		return
	}

	token, user, err := h.service.LoginWithGoogle(c.Request.Context(), req.IDToken)
	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login con Google exitoso",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
