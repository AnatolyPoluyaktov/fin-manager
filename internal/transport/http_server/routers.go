package http_server

import "github.com/go-chi/chi/v5"

func NewExpenseRouter(expenseHandler *ExpenseHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", expenseHandler.CreateExpense)
	return router
}

func NewCategoryRouter(categoryHandler *CategoryHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", categoryHandler.CreateCategory)
	router.Get("/", categoryHandler.GetCategories)
	return router
}
