package services

import "track-my-money/internal/core/domain/ports"

type service struct {
	extractor ports.DocumentExtractor
}

func NewService(ext ports.DocumentExtractor) ports.Service {
	return &service{
		extractor: ext,
	}
}
