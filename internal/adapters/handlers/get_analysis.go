package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetAnalysis(c *gin.Context) {

	analysisType := c.Query("type")
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	userUUID, _ := uuid.Parse("296f368f-f7b4-4388-8934-209e146de03c")

	switch analysisType {
	case "monthly-summary":
		report, err := h.service.GetMonthlyComparison(c.Request.Context(), userUUID, month)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, report)
	case "category-distribution": // <-- El nuevo type
		report, err := h.service.GetCategoryDistribution(c.Request.Context(), userUUID, month)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, report)
	default:
		c.JSON(400, gin.H{"error": "Tipo de análisis no válido o no especificado"})
	}
}
