package services

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"track-my-money/internal/core/domain"
)

func (s *service) ExtractStatementData(ctx context.Context, file io.Reader, fileName string, password string, bankIDStr string) ([]domain.Movement, string, error) {

	// 1. Convertir el ID del banco a entero
	bankID, err := strconv.Atoi(bankIDStr)
	if err != nil {
		return nil, "", fmt.Errorf("ID de banco inválido: %v", err)
	}

	// 2. Extraer el texto crudo del PDF
	text, err := s.extractor.Extract(file, password)
	if err != nil {
		return nil, "", err
	}

	// 3. Buscar si el banco ya tiene una regla aprobada en la base de datos
	rule, _ := s.repo.GetBankRule(ctx, bankID) // Si hay error o no existe, rule será nil

	// 4. Obtener el parser adecuado (Dinámico o IA)
	p := s.parserFactory.GetParser(rule)

	// 5. Procesar el texto
	movements, periodMonth, suggestedRule, err := p.Parse(text)
	if err != nil {
		return nil, "", err
	}

	// 6. ¡MAGIA!: Si el parser fue la IA, nos devolverá una regla sugerida. La guardamos.
	if suggestedRule != nil {
		suggestedRule.BankID = bankID
		suggestedRule.Status = "pending"
		// Guardamos en BD de forma asíncrona o síncrona
		_ = s.repo.SaveBankRule(ctx, suggestedRule)
	}

	// 7. Categorizar los movimientos
	for i := range movements {
		s.categorizer.Categorize(&movements[i])
	}

	return movements, periodMonth, nil
}
