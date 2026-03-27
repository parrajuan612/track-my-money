package ports

import "track-my-money/internal/core/domain"

type MovementParser interface {
	Parse(text string) ([]domain.Movement, string, error)
}
