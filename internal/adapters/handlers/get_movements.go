package handlers

import (
	"net/http"
	"strconv"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetMovements(c *gin.Context) {
	userUUID, _ := uuid.Parse("296f368f-f7b4-4388-8934-209e146de03c")

	filters := domain.MovementFilters{
		UserID: userUUID,
	}
	if page, err := strconv.Atoi(c.Query("page")); err == nil {
		filters.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil {
		filters.PageSize = pageSize
	}
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = &t
		}
	}

	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		if id, err := uuid.Parse(accountIDStr); err == nil {
			filters.AccountID = &id
		}
	}
	if bankIDStr := c.Query("bank_id"); bankIDStr != "" {
		if id, err := strconv.Atoi(bankIDStr); err == nil {
			filters.BankID = &id
		}
	}
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if id, err := strconv.Atoi(categoryIDStr); err == nil {
			filters.CategoryID = &id
		}
	}
	if movType := c.Query("type"); movType != "" {
		filters.Type = &movType
	}
	if query := c.Query("query"); query != "" {
		filters.Query = &query
	}
	filters.SortBy = c.Query("sort_by")
	filters.SortOrder = c.Query("sort_order")

	movements, total, err := h.service.GetMovements(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": movements,
		"meta": gin.H{
			"total":     total,
			"page":      filters.Page,
			"page_size": filters.PageSize,
		},
	})
}
