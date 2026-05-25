package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// Validar JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos. Revisa el correo y asegúrate de que la contraseña tenga al menos 6 caracteres."})
		return
	}

	// Llamar al servicio
	token, user, err := h.service.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		// Retornamos 409 Conflict si el correo ya existe, o 400 para otros errores
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Respuesta exitosa (201 Created)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Cuenta creada exitosamente",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
