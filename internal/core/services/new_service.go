package services

import "track-my-money/internal/core/domain/ports"

type service struct {
	extractor     ports.DocumentExtractor
	parserFactory ports.MovementParserFactory
	categorizer   ports.Categorizer
	repository    ports.Repository
	analysisRepo  ports.AnalysisRepository
}

func NewService(ext ports.DocumentExtractor, fact ports.MovementParserFactory, cat ports.Categorizer, rep ports.Repository, analy ports.AnalysisRepository) ports.Service {
	return &service{
		extractor:     ext,
		parserFactory: fact,
		categorizer:   cat,
		repository:    rep,
		analysisRepo:  analy,
	}
}
