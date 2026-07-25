package api

import (
	"context"
	"fmt"
	"os" // <-- Agregado para poder leer el puerto de Railway
	"track-my-money/internal/adapters/handlers"
	"track-my-money/internal/adapters/parsers"
	"track-my-money/internal/adapters/repository/postgres"

	"track-my-money/cmd/database"
	"track-my-money/internal/core/domain"
	"track-my-money/internal/core/services"

	"github.com/gin-gonic/gin"
)

func InitServer(r *gin.Engine) {
	// 1. Configuramos el puerto dinámico para Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "9080" // Puerto por defecto si lo corres en tu PC
	}

	r.LoadHTMLGlob("web/templates/*.html")

	db, err := database.InitPostgres()
	if err != nil {
		panic(fmt.Sprintf("No se pudo conectar a la DB:%v", err))
	}
	extractor := parsers.NewPDFExtractor()
	parserFactory := parsers.NewParserFactory()
	rules := domain.GetDefaultRules()
	categories := []domain.Category{
		{ID: 1, Name: "Restaurante/Comida/Supermercado"},
		{ID: 2, Name: "Ropa/Calzado"},
		{ID: 3, Name: "Salud/Deporte"},
		{ID: 4, Name: "Servicios/Suscripciones"},
		{ID: 5, Name: "Transporte"},
		{ID: 10, Name: "Salario"},
		{ID: 11, Name: "Otros Ingresos"},
		{ID: 12, Name: "Otros Gastos"},
	}
	categorizer := services.NewCategorizer(rules, categories)
	repo := postgres.NewPostgresRepository(db)

	service := services.NewService(extractor, parserFactory, categorizer, repo)
	handler := handlers.NewHandler(service)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.Use(func(c *gin.Context) {
		c.Set("ctx", ctx)
		c.Next()
	})

	InitRouter(r, handler)

	// 2. Usamos la variable port
	r.Run(fmt.Sprintf(":%s", port))
}