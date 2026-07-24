package handlers

import (
	"net/http"
	"strconv"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetAnalysis(c *gin.Context) {
	// 1. Extraemos el ID del usuario real desde el token JWT
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

	analysisType := c.Query("type")
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	switch analysisType {
	case "monthly-summary":
		report, err := h.service.GetMonthlyComparison(c.Request.Context(), userUUID, month)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, report)
	case "category-distribution":
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		report, err := h.service.GetCategoryDistribution(c.Request.Context(), userUUID, startDate, endDate)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, report)
	case "money-flow":
		timeRange := c.DefaultQuery("range", "1m")
		accountIDStr := c.Query("account_id")
		bankIDStr := c.Query("bank_id")
		var bankID *int
		var accountID *uuid.UUID

		if accountIDStr != "" {
			id, _ := uuid.Parse(accountIDStr)
			accountID = &id
		}
		if bankIDStr != "" {
			id, err := strconv.Atoi(bankIDStr)
			if err == nil {
				bankID = &id
			}
		}
		req := domain.MoneyFlowRequest{
			Range:     timeRange,
			AccountID: accountID,
			BankID:    bankID,
		}
		data, err := h.service.GetMoneyFlow(c.Request.Context(), userUUID, req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, data)
	case "trends":
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		var catID *int
		if catStr := c.Query("category_id"); catStr != "" {
			if id, err := strconv.Atoi(catStr); err == nil {
				catID = &id
			}
		}

		trends, err := h.service.GetExpenseTrends(c.Request.Context(), userUUID, startDate, endDate, catID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, trends)
	default:
		c.JSON(400, gin.H{"error": "Tipo de análisis no válido o no especificado"})
	}
}
