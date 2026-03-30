package postgres

import (
	"context"
	"track-my-money/internal/core/domain"
	"track-my-money/internal/core/domain/ports"

	"gorm.io/gorm"
)

type statementRepository struct {
	db *gorm.DB
}

func NewStatementRepository(db *gorm.DB) ports.Repository {
	return &statementRepository{db: db}
}
func (r *statementRepository) SaveStatementWithMovements(ctx context.Context, stmt *domain.Statement, movs []domain.Movement) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		stmtModel := StatementModel{
			UserID:      stmt.UserID,
			AccountID:   stmt.AccountID,
			BankID:      stmt.BankID,
			FileName:    stmt.FileName,
			PeriodMonth: stmt.PeriodMonth,
			Status:      string(domain.StatusPending),
		}
		if err := tx.Create(&stmtModel).Error; err != nil {
			return err
		}

		movModels := make([]MovementModel, len(movs))
		for i, m := range movs {
			movModels[i] = MovementModel{
				UserID:      m.UserID,
				AccountID:   m.AccountID,
				StatementID: &stmtModel.ID,
				CategoryID:  m.CategoryID,
				Date:        m.Date,
				Description: m.Description,
				Amount:      m.Amount,
				Type:        string(m.Type),
			}
		}
		if len(movModels) > 0 {
			if err := tx.Create(&movModels).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
