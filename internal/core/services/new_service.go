package services

import "track-my-money/internal/core/domain/ports"

type service struct {
	extractor     ports.DocumentExtractor
	parserFactory ports.MovementParserFactory
	categorizer   ports.Categorizer 
}

func NewService(ext ports.DocumentExtractor, fact ports.MovementParserFactory, cat ports.Categorizer) ports.Service {
	return &service{
		extractor:     ext,
		parserFactory: fact,
		categorizer:   cat,
	}
}
