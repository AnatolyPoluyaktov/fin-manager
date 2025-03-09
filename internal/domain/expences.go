package domain

// Перечисление валют
type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	RUB Currency = "RUB"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

type NewExpenseData struct {
	CategoryID int
	Currency   Currency
	Amount     int
	Note       string
	ActionDate string
}

// Структура расходов
type Expense struct {
	ID         int
	CategoryID int
	Currency   Currency
	Amount     int
	Note       string
	ActionDate string
}

type ExpenseWithCategory struct {
	ID           int
	CategoryName string
	Currency     Currency
	Amount       int
	Note         string
	ActionDate   string
}

func NewExpense(data NewExpenseData) (Expense, error) {
	return Expense{
		CategoryID: data.CategoryID,
		Currency:   data.Currency,
		Amount:     data.Amount,
		Note:       data.Note,
		ActionDate: data.ActionDate,
	}, nil
}
