package api

import (
	"context"
	"fmt"
	"track-my-money/internal/adapters/handlers"
	"track-my-money/internal/adapters/parsers"
	"track-my-money/internal/adapters/repository/postgres"

	"track-my-money/cmd/database"
	"track-my-money/internal/core/domain"
	"track-my-money/internal/core/services"

	"github.com/gin-gonic/gin"
)

var (
	port = ("9080")
)

func InitServer(r *gin.Engine) {

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
	repo := postgres.NewStatementRepository(db)
	analysisRepo := postgres.NewAnalysisRepository(db)
	service := services.NewService(extractor, parserFactory, categorizer, repo, analysisRepo)
	handler := handlers.NewHandler(service)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.Use(func(c *gin.Context) {
		c.Set("ctx", ctx)
		c.Next()
	})

	InitRouter(r, handler)

	r.Run(fmt.Sprintf(":%s", port))

}
