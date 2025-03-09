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

const testToken = "my-secret-token"

func setupCategoryRouter() *chi.Mux {
	repo, _ := repositories.NewCategoryRepository(testDB.DB)
	usecase := usecase.NewCategoryUseCase(repo)
	handler := http_server.NewCategoryHandler(usecase)
	router := chi.NewRouter()
	router.Use(middleware.AuthMiddleware(TestToken)) // Добавляем middleware аутентификации
	router.Post("/api/v1/categories", handler.CreateCategory)
	router.Get("/api/v1/categories", handler.GetCategories)
	return router
}

// ✅ Тест создания категории с токеном
func TestCreateCategory(t *testing.T) {
	router := setupCategoryRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	payload := `{"name": "test_category", "description": "test_description"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/categories", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken) // Добавляем токен

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}

	var createCategoryResponse http_server.CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&createCategoryResponse)
	if err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "test_category", createCategoryResponse.Name)
	assert.Equal(t, "test_description", createCategoryResponse.Description)
	assert.NotEmpty(t, createCategoryResponse.ID)
}

// ✅ Тест получения списка категорий с токеном
func TestGetCategories(t *testing.T) {
	expectedCategories := []http_server.CategoryResponse{
		{ID: 1, Name: "test_category", Description: "test_description"},
		{ID: 2, Name: "test_category2", Description: "test_description2"},
	}

	// Вставляем тестовые данные
	for _, category := range expectedCategories {
		testDB.DB.Pool.Exec(context.Background(), "INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)", category.ID, category.Name, category.Description)
	}

	router := setupCategoryRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/api/v1/categories", nil)
	req.Header.Set("Authorization", "Bearer "+testToken) // Добавляем токен

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}

	var getCategoriesResponse []http_server.CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&getCategoriesResponse)
	if err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, getCategoriesResponse)
	assert.ElementsMatch(t, expectedCategories, getCategoriesResponse)
}

// ✅ Тест запроса без токена (должен вернуть 401 Unauthorized)
func TestUnauthorizedAccess(t *testing.T) {
	router := setupCategoryRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/api/v1/categories", nil) // Без токена

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode) // Ожидаем 401 Unauthorized
}
