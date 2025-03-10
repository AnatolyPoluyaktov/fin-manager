package repositories

import (
	"context"
	"fin-manager/internal/domain"
	"fin-manager/internal/storage/pg"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type ExpensesRepository struct {
	DB *pg.DB
}

func NewExpensesRepository(db *pg.DB) (*ExpensesRepository, error) {
	return &ExpensesRepository{DB: db}, nil
}

func (e *ExpensesRepository) CreateExpense(ctx context.Context, new_expense domain.NewExpenseData) (*domain.Expense, error) {
	tx, err := e.DB.Pool.Begin(ctx)
	if err != nil {
		_ = fmt.Errorf("error starting transaction: %v", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

			_ = fmt.Errorf("error rolling back transaction: %v", err)
		}
	}(tx, ctx)

	created_expense, err := domain.NewExpense(new_expense)
	if err != nil {
		return nil, fmt.Errorf("error creating new expense: %v", err)
	}
	query := `INSERT INTO expenses (amount, currency, category_id, action_date, note) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = tx.QueryRow(ctx, query, new_expense.Amount, new_expense.Currency, new_expense.CategoryID, new_expense.ActionDate, new_expense.Note).Scan(&created_expense.ID)
	if err != nil {
		return nil, fmt.Errorf("error inserting expense: %v", err)
	}

	// Фиксируем транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &created_expense, nil
}
