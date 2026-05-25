package api

import (
	"track-my-money/internal/adapters/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, handler *handlers.Handler) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.GET("/", handler.Home)
	v1 := r.Group("/api/v1")
	{

		auth := v1.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/google", handler.LoginGoogle)
			auth.POST("/register", handler.Register)
		}

		v1.POST("/statements/parse", handler.ExtractStatementData)
		v1.POST("/movements/bulk", handler.SaveMovements)
		v1.GET("/analysis", handler.GetAnalysis)
		v1.GET("/movements", handler.GetMovements)
	}
}
