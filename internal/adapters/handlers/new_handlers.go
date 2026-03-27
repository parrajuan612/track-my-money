package handlers

import (
	"net/http"
	"track-my-money/internal/core/domain/ports"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ports.Service
}

func NewHandler(service ports.Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
