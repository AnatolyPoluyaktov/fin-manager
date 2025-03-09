package usecase

import (
	"context"
	"fin-manager/internal/domain"
)

type ExpenseRepository interface {
	CreateExpense(ctx context.Context, new_expense domain.NewExpenseData) (*domain.Expense, error)
}

type ExpenseUseCase struct {
	ExpenseRepository ExpenseRepository
}

func NewExpenseUseCase(repo ExpenseRepository) *ExpenseUseCase {

	return &ExpenseUseCase{ExpenseRepository: repo}
}

func (e *ExpenseUseCase) CreateExpense(ctx context.Context, new_expense domain.NewExpenseData) (*domain.Expense, error) {
	return e.ExpenseRepository.CreateExpense(ctx, new_expense)
}
