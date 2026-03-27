package services

import (
	"fmt"
	"io"
	"track-my-money/internal/core/domain"
)

func (s *service) ExtractStatementData(file io.Reader, fileName string, password string, bankID string) ([]domain.Movement, string, error) {

	text, err := s.extractor.Extract(file, password)
	if err != nil {
		return nil, "", err
	}
	fmt.Println(text)
	return []domain.Movement{}, "", nil
}
