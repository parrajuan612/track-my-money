package services

import (
	"track-my-money/internal/core/domain/ports"
)

type service struct {
	extractor     ports.DocumentExtractor
	parserFactory ports.MovementParserFactory
	categorizer   ports.Categorizer
	repo          ports.AppRepository
}

func NewService(
	extractor ports.DocumentExtractor,
	parserFactory ports.MovementParserFactory,
	categorizer ports.Categorizer,
	repo ports.AppRepository,
) ports.Service {
	return &service{
		extractor:     extractor,
		parserFactory: parserFactory,
		categorizer:   categorizer,
		repo:          repo,
	}
}
