package api

import (
	"context"
	"fmt"
	"track-my-money/internal/adapters/handlers"
	"track-my-money/internal/adapters/parsers"

	"track-my-money/internal/core/domain"
	"track-my-money/internal/core/services"

	"github.com/gin-gonic/gin"
)

var (
	port = ("9080")
)

func InitServer(r *gin.Engine) {

	r.LoadHTMLGlob("web/templates/*.html")
	extractor := parsers.NewPDFExtractor()
	parserFactory := parsers.NewParserFactory()
	rules := domain.GetDefaultRules()
	var categories []domain.Category
	categorizer := services.NewCategorizer(rules, categories)
	service := services.NewService(extractor, parserFactory, categorizer)
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
