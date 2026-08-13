package postgres

import (
	"context"
	"track-my-money/internal/core/domain"
)

func (r *postgresRepository) GetBankRule(ctx context.Context, bankID int) (*domain.BankParsingRule, error) {
	var rule domain.BankParsingRule
	err := r.db.WithContext(ctx).Where("bank_id = ? AND status = ?", bankID, "approved").First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *postgresRepository) SaveBankRule(ctx context.Context, rule *domain.BankParsingRule) error {
	// Guarda la regla nueva o la actualiza si ya existía para ese banco
	return r.db.WithContext(ctx).Save(rule).Error
}
