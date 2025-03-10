package tests

import (
	"context"
	"encoding/json"
	"fin-manager/internal/repositories"
	"fin-manager/internal/transport/http_server"
	"fin-manager/internal/transport/http_server/middleware"
	"fin-manager/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupRouter() *chi.Mux {
	repo, _ := repositories.NewExpensesRepository(testDB.DB)
	useCase := usecase.NewExpenseUseCase(repo)
	handler := http_server.NewExpenseHandler(useCase)

	router := chi.NewRouter()
	router.Use(middleware.AuthMiddleware(TestToken))
	router.Post("/api/v1/expenses", handler.CreateExpense)
	return router
}

func TestCreateExpense(t *testing.T) {
	_, err2 := testDB.DB.Pool.Exec(context.Background(), "INSERT INTO categories (id, name, description) VALUES (1, 'test_category', 'test_description')")
	if err2 != nil {
		return
	}

	router := setupRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	payload := `{"amount": 100, "currency": "RUB", "category_id": 1, "action_date": "2022-01-01", "note": "test"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/expenses", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken) // Добавляем токен

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}

	var createExpenseResponse http_server.ExpenseResponse
	err = json.NewDecoder(resp.Body).Decode(&createExpenseResponse)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, 100, createExpenseResponse.Amount)
	assert.Equal(t, "RUB", createExpenseResponse.Currency)
	assert.Equal(t, 1, createExpenseResponse.CategoryID)
	assert.Equal(t, "2022-01-01", createExpenseResponse.ActionDate)
	assert.Equal(t, "test", createExpenseResponse.Note)
	assert.NotEmpty(t, createExpenseResponse.ID)
}

func TestUnauthorizedExpenseCreation(t *testing.T) {
	router := setupRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	payload := `{"amount": 100, "currency": "RUB", "category_id": 1, "action_date": "2022-01-01", "note": "test"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/expenses", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	// ⛔ Не передаем токен!

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode) // Ожидаем 401 Unauthorized
}
