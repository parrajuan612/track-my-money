package domain

import (
	"time"

	"github.com/google/uuid"
)

type BankParsingRule struct {
	ID             uuid.UUID `json:"id"`
	BankID         int       `json:"bank_id"`
	RegexRow       string    `json:"regex_row"`
	RegexDate      string    `json:"regex_date"`
	RegexAmount    string    `json:"regex_amount"`
	IncomeKeywords string    `json:"income_keywords"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
