package handlers

import (
	"net/http"
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/banks
func (h *Handler) GetBanks(c *gin.Context) {
	banks, err := h.service.GetBanks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la lista de bancos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": banks})
}

// GET /api/v1/accounts
func (h *Handler) GetAccounts(c *gin.Context) {
	// Leemos el user_id que colocó el AuthMiddleware
	userID, _ := c.Get("user_id")

	accounts, err := h.service.GetAccountsByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las cuentas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": accounts})
}

// POST /api/v1/accounts
func (h *Handler) CreateAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req domain.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos para crear la cuenta"})
		return
	}

	newAccount := &domain.Account{
		UserID: userID.(string),
		BankID: req.BankID,
		Name:   req.Name,
	}

	if err := h.service.CreateAccount(c.Request.Context(), newAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la cuenta en la base de datos"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Cuenta creada exitosamente", "data": newAccount})
}
