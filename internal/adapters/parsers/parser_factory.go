package parsers

import (
	"fmt"
	"track-my-money/internal/core/domain/ports"
)

type parserFactory struct{}

func NewParserFactory() ports.MovementParserFactory {
	return &parserFactory{}
}

const (
	BankBancolombia = "1"
	BankNu          = "2"
)

func (f *parserFactory) GetParser(bankID string) (ports.MovementParser, error) {

	switch bankID {

	case BankBancolombia:
		return NewBancolombiaParser(), nil

	case BankNu:
		return NewNuParser(), nil

	default:
		return nil, fmt.Errorf("unsupported bank")
	}
}
