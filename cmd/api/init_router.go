package api

import (
	"track-my-money/internal/adapters/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, handler *handlers.Handler) {
	// Configuración de CORS
r.Use(cors.New(cors.Config{
		// Agregamos tu URL de Vercel y dejamos localhost por si quieres seguir probando en tu PC
		AllowOrigins:     []string{"http://localhost:3000", "https://track-my-money-web.vercel.app"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
	}))
	r.GET("/", handler.Home)

	v1 := r.Group("/api/v1")
	{
		// 🟢 RUTAS PÚBLICAS (No necesitan Token)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/google", handler.LoginGoogle)
			auth.POST("/register", handler.Register)
		}

		// 🔴 RUTAS PRIVADAS (Protegidas por el Portero)
		protected := v1.Group("/")
		protected.Use(handlers.AuthMiddleware()) // Aquí actúa nuestro "Portero"
		{
			// Extractos
			protected.POST("/statements/parse", handler.ExtractStatementData)
			protected.POST("/movements/bulk", handler.SaveMovements)

			// Dashboard y Análisis
			protected.GET("/analysis", handler.GetAnalysis)

			// Movimientos
			protected.GET("/movements", handler.GetMovements)
			protected.PUT("/movements/:id", handler.UpdateMovement)
			protected.POST("/movements", handler.CreateMovement)
			protected.DELETE("/movements/:id", handler.DeleteMovement)
			// Bancos y Cuentas
			protected.GET("/banks", handler.GetBanks)
			protected.GET("/accounts", handler.GetAccounts)
			protected.POST("/accounts", handler.CreateAccount)

			//Session
			protected.GET("/auth/me", handler.GetMe)
		}
	}
}
