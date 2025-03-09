package http_server

import "fin-manager/internal/domain"

func toDomainNewExpense(expenseRequest ExpenseRequest) (domain.NewExpenseData, error) {
	return domain.NewExpenseData{
		Amount:     expenseRequest.Amount,
		CategoryID: expenseRequest.CategoryID,
		Currency:   domain.Currency(expenseRequest.Currency),
		ActionDate: expenseRequest.ActionDate,
		Note:       expenseRequest.Note,
	}, nil
}

func toResponseExpense(expense domain.Expense) ExpenseResponse {
	return ExpenseResponse{
		ID:         expense.ID,
		Amount:     expense.Amount,
		CategoryID: expense.CategoryID,
		Currency:   string(expense.Currency),
		ActionDate: expense.ActionDate,
		Note:       expense.Note,
	}
}

func toDomainNewCategory(categoryRequest CategoryRequest) (domain.NewCategoryData, error) {
	return domain.NewCategoryData{
		Name:        categoryRequest.Name,
		Description: categoryRequest.Description,
	}, nil
}

func toResponseCategory(category domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}

func toResponseCategories(categories []domain.Category) []CategoryResponse {
	response := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
	return response
}
