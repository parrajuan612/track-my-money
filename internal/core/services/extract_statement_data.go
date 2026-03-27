package services

import (
	"io"
	"track-my-money/internal/core/domain"
)

func (s *service) ExtractStatementData(file io.Reader, fileName string, password string, bankID string) ([]domain.Movement, string, error) {

	text, err := s.extractor.Extract(file, password)
	if err != nil {
		return nil, "", err
	}
	p, err := s.parserFactory.GetParser(bankID)
	if err != nil {
		return nil, "", err
	}
	movements, periodMonth, err := p.Parse(text)

	for i := range movements {
		s.categorizer.Categorize(&movements[i])
	}

	return movements, periodMonth, nil

}
