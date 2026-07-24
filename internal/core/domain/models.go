package domain

type SaveStatementRequest struct {
	AccountID   string               `json:"account_id" binding:"required"`
	BankID      int                  `json:"bank_id" binding:"required"`
	PeriodMonth string               `json:"period_month" binding:"required"`
	FileName    string               `json:"file_name" binding:"required"`
	Movements   []MovementRequestDTO `json:"movements" binding:"required,dive"`
}

type MovementRequestDTO struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
}
