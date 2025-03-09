package http_server

import (
	"context"
	"encoding/json"
	"fin-manager/internal/common/server"
	"fin-manager/internal/domain"
	"net/http"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, new_category domain.NewCategoryData) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
}

type CategoryHandler struct {
	CategoryUseCase CategoryUseCase
}

func NewCategoryHandler(useCase CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{CategoryUseCase: useCase}
}

func (c *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryRequest CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&categoryRequest); err != nil {
		server.BadRequest("invalid_json", err, w, r)
	}
	if err := categoryRequest.Validate(); err != nil {
		server.BadRequest("invalid_request", err, w, r)
		return
	}
	new_category_data, err := toDomainNewCategory(categoryRequest)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}
	created_category, err := c.CategoryUseCase.CreateCategory(r.Context(), new_category_data)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}
	response := toResponseCategory(*created_category)
	server.RespondOK(response, w, r)

}

func (c *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.CategoryUseCase.GetCategories(r.Context())

	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}
	response := toResponseCategories(categories)

	server.RespondOK(response, w, r)
}
