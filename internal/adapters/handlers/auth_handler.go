func (h *Handler) GetMe(c *gin.Context) {
	// 1. El portero (middleware) ya validó el token y nos dejó el ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	// 2. Convertimos el ID (si viene como string desde el token) a UUID según sea necesario.
	// Asumiendo que lo manejas como string y tienes un método GetUserByID en tu servicio:
	user, err := h.service.GetUserByID(c.Request.Context(), userIDStr.(string))
	
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado en la base de datos"})
		return
	}

	// 3. Devolvemos los datos REALES extraídos de la base de datos
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,  // ¡AQUÍ VIAJA EL NOMBRE REAL (ej: "juan camilo...")!
			"email": user.Email,
		},
	})
}
