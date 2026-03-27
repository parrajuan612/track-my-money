package ports

import "track-my-money/internal/core/domain"

type Categorizer interface {
	Categorize(mov *domain.Movement)
}
