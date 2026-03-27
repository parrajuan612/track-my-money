package api

import (
	"track-my-money/internal/adapters/handlers"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, handler *handlers.Handler) {
	r.GET("/", handler.Home)
	v1 := r.Group("/api/v1")
	{
		v1.POST("/statements/parse", handler.ExtractStatementData)
	}
}
