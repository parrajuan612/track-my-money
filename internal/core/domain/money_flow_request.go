package domain

import "github.com/google/uuid"

type MoneyFlowRequest struct {
	Range     string // "1w", "1m", "1y"
	AccountID *uuid.UUID
	BankID    *int
}
