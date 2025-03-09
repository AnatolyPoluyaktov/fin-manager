package http_server

import (
	"context"
	"encoding/json"
	"fin-manager/internal/common/server"
	"fin-manager/internal/domain"
	"net/http"
)

type ExpenseUseCase interface {
	CreateExpense(ctx context.Context, new_expense domain.NewExpenseData) (*domain.Expense, error)
}

type ExpenseHandler struct {
	ExpenseUseCase ExpenseUseCase
}

func NewExpenseHandler(useCase ExpenseUseCase) *ExpenseHandler {
	return &ExpenseHandler{ExpenseUseCase: useCase}
}

func (e *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expenseRequest ExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&expenseRequest); err != nil {
		server.BadRequest("invalid_json", err, w, r)
	}

	if err := expenseRequest.Validate(); err != nil {
		server.BadRequest("invalid_request", err, w, r)
		return
	}

	new_expense_data, err := toDomainNewExpense(expenseRequest)

	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	created_expense, err := e.ExpenseUseCase.CreateExpense(r.Context(), new_expense_data)

	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	response := toResponseExpense(*created_expense)

	server.RespondOK(response, w, r)

}
