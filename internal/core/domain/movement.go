package domain

import (
	"time"

	"github.com/google/uuid"
)

type Movement struct {
	ID           uuid.UUID    `json:"id"`
	UserID       uuid.UUID    `json:"user_id"`
	AccountID    uuid.UUID    `json:"account_id"`
	StatementID  *uuid.UUID   `json:"statement_id,omitempty"`
	CategoryID   int          `json:"category_id"`
	CategoryName string       `json:"category_name"`
	Date         time.Time    `json:"date"`
	Description  string       `json:"description"`
	Amount       float64      `json:"amount"`
	Type         MovementType `json:"type"`
	CreatedAt    time.Time    `json:"created_at"`
}
type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}
type MovementType string

const (
	TypeIncome  MovementType = "income"
	TypeExpense MovementType = "expense"
)

type MovementFilters struct {
	UserID     uuid.UUID
	StartDate  *time.Time
	EndDate    *time.Time
	BankID     *int
	AccountID  *uuid.UUID
	CategoryID *int
	Type       *string
	Query      *string
	Page       int
	PageSize   int
	SortBy     string
	SortOrder  string
}

type MovementTable struct {
	ID          uuid.UUID
	Date        time.Time
	Description string
	Amount      float64
	Type        string
	// Campos legibles para el usuario
	CategoryName string `json:"category_name"`
	AccountName  string `json:"account_name"`
	BankName     string `json:"bank_name"`
}

type UpdateMovementRequest struct {
	Description string  `json:"description" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	CategoryID  int     `json:"category_id" binding:"required"`
}

type CreateMovementRequest struct {
	Date        string  `json:"date" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	CategoryID  int     `json:"category_id" binding:"required"`
}
