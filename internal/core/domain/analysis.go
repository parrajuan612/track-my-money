package domain

type MonthTotals struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type ComparisonResponse struct {
	CurrentMonth  MonthTotals `json:"current_month"`
	PreviousMonth MonthTotals `json:"previous_month"`
	Variations    Variations  `json:"variations"`
}

type Variations struct {
	IncomePercentage  float64 `json:"income_percentage"`
	ExpensePercentage float64 `json:"expense_percentage"`
}

type CategoryReport struct {
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

type CategoryDistributionResponse struct {
	TotalExpenses float64          `json:"total_expenses"`
	Categories    []CategoryReport `json:"categories"`
}
type TimeSeriesData struct {
	Label   string  `json:"label"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type MoneyFlowResponse struct {
	PeriodType string           `json:"period_type"`
	Data       []TimeSeriesData `json:"data"`
}
type TrendReport struct {
	Name   string  `gorm:"column:name" json:"name"`
	Amount float64 `gorm:"column:amount" json:"amount"`
}
