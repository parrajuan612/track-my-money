package parsers

import (
	"track-my-money/internal/core/domain"
	"track-my-money/internal/core/domain/ports"
)

type parserFactory struct{}

func NewParserFactory() ports.MovementParserFactory {
	return &parserFactory{}
}

// Ahora solo recibimos la regla (que el servicio ya buscó en la BD)
func (f *parserFactory) GetParser(rule *domain.BankParsingRule) ports.MovementParser {

	// 1. Si existe una regla para este banco y está aprobada, usamos el motor súper rápido
	if rule != nil && rule.Status == "approved" {
		return NewDynamicParser(*rule)
	}

	// 2. Si el banco es nuevo y no tiene regla (o está pendiente), usamos la IA
	return NewOpenAIParser()
}
