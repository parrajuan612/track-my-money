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

	movements, periodMonth, err := h.service.ExtractStatementData(file, header.Filename, password, bankID)
	if err != nil {

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"Movements":   movements,
		"Periodmonth": periodMonth,
		"Filename":    header.Filename,
	})
}
