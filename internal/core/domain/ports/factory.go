package ports

import "track-my-money/internal/core/domain"

type MovementParserFactory interface {
	GetParser(rule *domain.BankParsingRule) MovementParser
}
