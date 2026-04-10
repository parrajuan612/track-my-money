package postgres

import (
	"track-my-money/internal/core/domain/ports"

	"gorm.io/gorm"
)

// Este es el único struct que usaremos para todo
type postgresRepository struct {
	db *gorm.DB
}

// Este es el único constructor que llamarás desde InitServer
func NewPostgresRepository(db *gorm.DB) ports.AppRepository {
	return &postgresRepository{db: db}
}
