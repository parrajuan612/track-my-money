package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ExtractStatementData(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "archivo requerido")
		return
	}
	defer file.Close()

	password := c.PostForm("password")
	bankID := c.PostForm("bank_id")

	// Pasamos c.Request.Context() como primer parámetro
	movements, periodMonth, err := h.service.ExtractStatementData(c.Request.Context(), file, header.Filename, password, bankID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         movements,
		"period_month": periodMonth,
		"file_name":    header.Filename,
	})
}
